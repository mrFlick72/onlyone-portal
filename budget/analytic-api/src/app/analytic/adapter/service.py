from typing import List

from app.analytic.domain.expense import ExpenseRecord
from app.analytic.domain.service import ExpenseLoader
import httpx
from datetime import date
from app.user.domain.user_name_resolver import UserNameResolver

class RestExpenseLoader(ExpenseLoader):
    
    def __init__(self, budget_api_base_url: str, access_token_resolver: UserNameResolver) -> None:
        self.budget_api_base_url = budget_api_base_url
        self.access_token_resolver = access_token_resolver

    def expenseFor(self, year: int, month: int, tags: List[str]) -> List[ExpenseRecord]:
        response = httpx.put(
            f"{self.budget_api_base_url}/api/budget/expense",
            json={"month": str(month), "year": str(year), "searchTagList": tags},
            headers={"Authorization": f"Bearer {self.access_token_resolver}"},
        )
        response.raise_for_status()
        
        records: List[ExpenseRecord] = []
        for daily in response.json().get("dailyBudgetExpenseRepresentationList", []):
            for item in daily.get("budgetExpenseRepresentationList", []):
                records.append(ExpenseRecord(
                    id=item["id"],
                    date=date.fromisoformat(item["date"]),
                    amount=float(item["amount"]),
                    note=item["note"],
                    tag_values=[t["searchTagValue"] for t in item.get("tags", [])],
                ))
        return records
