from concurrent.futures import ThreadPoolExecutor
from typing import Any, Dict, List, Optional

import httpx

from app.analytic.domain.expense import ExpenseTag, ProjectedExpense
from app.analytic.domain.reindex import BudgetExpenseSource
from app.infrastructure.security.security_context_resolver import (
    SecurityContextResolver,
)
from app.money.domain.money import Money
from app.time.domain.date import Date

MAX_CONCURRENT_BUDGET_API_CALLS = 6
MONTHS = range(1, 13)


class RestBudgetExpenseSource(BudgetExpenseSource):
    """Pulls the current user's expenses from budget-api, propagating their JWT.

    `budget-api`'s expense endpoint is queried per (year, month); the caller's
    token and user name are captured once on the request thread (the security
    context lives in a thread-local) before fanning out to the pool.

    `transport` is a seam for tests; production passes none (plain httpx)."""

    def __init__(
        self,
        budget_api_base_url: str,
        security_context_resolver: SecurityContextResolver,
        transport: Optional[httpx.BaseTransport] = None,
    ) -> None:
        self.budget_api_base_url = budget_api_base_url
        self.security_context_resolver = security_context_resolver
        self._transport = transport

    def fetch_expenses(self, from_year: int, to_year: int) -> List[ProjectedExpense]:
        context = self.security_context_resolver.get_security_context()
        token = context.token
        user_name = context.user_name.content

        periods = [(year, month) for year in range(from_year, to_year + 1) for month in MONTHS]
        with httpx.Client(transport=self._transport) as client, ThreadPoolExecutor(
            max_workers=MAX_CONCURRENT_BUDGET_API_CALLS
        ) as pool:
            batches = list(
                pool.map(
                    lambda period: self._fetch(client, token, user_name, period[0], period[1]),
                    periods,
                )
            )
        return [expense for batch in batches for expense in batch]

    def _fetch(
        self,
        client: httpx.Client,
        token: str,
        user_name: str,
        year: int,
        month: int,
    ) -> List[ProjectedExpense]:
        response = client.put(
            f"{self.budget_api_base_url}/api/budget/expense",
            json={"month": str(month), "year": str(year), "searchTagList": []},
            headers={"Authorization": f"Bearer {token}"},
        )
        response.raise_for_status()
        return self._to_expenses(user_name, response.json())

    @staticmethod
    def _to_expenses(user_name: str, payload: Dict[str, Any]) -> List[ProjectedExpense]:
        expenses: List[ProjectedExpense] = []
        for daily in payload.get("dailyBudgetExpenseRepresentationList", []):
            for item in daily.get("budgetExpenseRepresentationList", []):
                expenses.append(
                    ProjectedExpense(
                        id=item["id"],
                        user_name=user_name,
                        date=Date.date_for(item["date"]),
                        amount=Money.money_for(item["amount"]),
                        note=item.get("note") or "",
                        # budget-api's REST layer serializes tags as
                        # {tagKey, tagValue} (web/tags SearchTagRepresentation),
                        # unlike the Kafka event which uses {key, value}.
                        tags=[
                            ExpenseTag(key=tag["tagKey"], value=tag["tagValue"])
                            for tag in item.get("tags", [])
                        ],
                    )
                )
        return expenses
