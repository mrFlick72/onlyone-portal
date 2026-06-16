from abc import ABC, abstractmethod
from typing import List

from app.analytic.domain.expense import ProjectedExpense
from app.analytic.domain.repository import ExpenseProjectionRepository


class BudgetExpenseSource(ABC):
    """Port over the authoritative expense data owned by budget-api, used to
    rebuild the projection. Implementations scope the fetch to the current user."""

    @abstractmethod
    def fetch_expenses(self, from_year: int, to_year: int) -> List[ProjectedExpense]:
        ...


class ExpenseReindexService:
    """Rebuilds the current user's projection from budget-api over a year range.

    Upsert-only: it fills gaps left by missed CREATE/UPDATE events but does not
    remove rows whose DELETE event was missed (those linger until re-touched)."""

    def __init__(
        self,
        source: BudgetExpenseSource,
        repository: ExpenseProjectionRepository,
    ) -> None:
        self.source = source
        self.repository = repository

    def reindex(self, from_year: int, to_year: int) -> int:
        expenses = self.source.fetch_expenses(from_year, to_year)
        for expense in expenses:
            self.repository.save(expense)
        return len(expenses)
