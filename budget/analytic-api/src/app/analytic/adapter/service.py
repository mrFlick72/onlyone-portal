from concurrent.futures import ThreadPoolExecutor
from typing import List

from app.analytic.domain.expense import BudgetExpenseAnalysisRequest, ExpenseRecord
from app.analytic.domain.service import ExpenseLoader
from app.infrastructure.security.security_context_resolver import (
    SecurityContextResolver,
)
from app.money.domain.money import Money
from app.time.domain.date import Date
import httpx

MAX_CONCURRENT_BUDGET_API_CALLS = 6


class RestExpenseLoader(ExpenseLoader):

    def __init__(
        self,
        budget_api_base_url: str,
        security_context_resolver: SecurityContextResolver,
    ) -> None:
        self.budget_api_base_url = budget_api_base_url
        self.security_context_resolver = security_context_resolver

    def expenseFor(self, request: BudgetExpenseAnalysisRequest) -> List[ExpenseRecord]:
        token = self.security_context_resolver.get_security_context().token
        with httpx.Client() as client:
            return self._fetch(client, token, request)

    def expenses_for_all(
        self, requests: List[BudgetExpenseAnalysisRequest]
    ) -> List[ExpenseRecord]:
        # the SecurityContext lives in a thread-local: capture the token here,
        # on the caller's thread, before fanning out to the pool
        token = self.security_context_resolver.get_security_context().token
        with httpx.Client() as client, ThreadPoolExecutor(
            max_workers=MAX_CONCURRENT_BUDGET_API_CALLS
        ) as pool:
            batches = list(
                pool.map(lambda request: self._fetch(client, token, request), requests)
            )
        return [record for batch in batches for record in batch]

    def _fetch(
        self,
        client: httpx.Client,
        token: str,
        request: BudgetExpenseAnalysisRequest,
    ) -> List[ExpenseRecord]:
        response = client.put(
            f"{self.budget_api_base_url}/api/budget/expense",
            json={
                "month": str(request.month),
                "year": str(request.year),
                "searchTagList": request.tags,
            },
            headers={"Authorization": f"Bearer {token}"},
        )
        response.raise_for_status()

        records: List[ExpenseRecord] = []
        for daily in response.json().get("dailyBudgetExpenseRepresentationList", []):
            for item in daily.get("budgetExpenseRepresentationList", []):
                records.append(
                    ExpenseRecord(
                        id=item["id"],
                        date=Date.date_for(item["date"]),
                        amount=Money.money_for(item["amount"]),
                        note=item["note"],
                        tag_values=[t["tagValue"] for t in item.get("tags", [])],
                    )
                )
        return records
