from typing import List

from app.analytic.domain.expense import ExpenseRecord
from app.analytic.domain.service import BudgetExpenseAnalysisRequest, ExpenseLoader
from app.infrastructure.security.security_context_resolver import (
    SecurityContextResolver,
)
import httpx
from datetime import date


class RestExpenseLoader(ExpenseLoader):

    def __init__(
        self,
        budget_api_base_url: str,
        security_context_resolver: SecurityContextResolver,
    ) -> None:
        self.budget_api_base_url = budget_api_base_url
        self.security_context_resolver = security_context_resolver

    def expenseFor(self, request: BudgetExpenseAnalysisRequest) -> List[ExpenseRecord]:
        response = httpx.put(
            f"{self.budget_api_base_url}/api/budget/expense",
            json={
                "month": str(request.month),
                "year": str(request.year),
                "searchTagList": request.tags,
            },
            headers={
                "Authorization": f"Bearer {self.security_context_resolver.get_security_context().token}"
            },
        )
        response.raise_for_status()

        records: List[ExpenseRecord] = []
        for daily in response.json().get("dailyBudgetExpenseRepresentationList", []):
            for item in daily.get("budgetExpenseRepresentationList", []):
                records.append(
                    ExpenseRecord(
                        id=item["id"],
                        date=date.fromisoformat(item["date"]),
                        amount=float(item["amount"]),
                        note=item["note"],
                        tag_values=[t["searchTagValue"] for t in item.get("tags", [])],
                    )
                )
        return records
