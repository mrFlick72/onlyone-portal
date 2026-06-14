from abc import ABC, abstractmethod
from typing import List, Optional

from app.analytic.domain.expense import ProjectedExpense, TagTotal, YearTotal


class ExpenseProjectionRepository(ABC):
    """Port over the read-optimised expense projection.

    Writes (`save`/`delete`) are driven by the Kafka consumer and carry the owning
    user in the event payload. Reads are scoped to a `user_name` resolved from the
    caller's JWT.
    """

    @abstractmethod
    def save(self, expense: ProjectedExpense) -> None:
        """Insert or replace the projected expense identified by `expense.id`."""
        ...

    @abstractmethod
    def delete(self, expense_id: str) -> None:
        """Remove the projected expense; a no-op when it is not present."""
        ...

    @abstractmethod
    def total_by_tag(
        self, user_name: str, year: int, month: Optional[int], tag_keys: List[str]
    ) -> List[TagTotal]:
        """Totals grouped by tag value for one year (optionally one month),
        restricted to expenses carrying any of `tag_keys` when it is non-empty."""
        ...

    @abstractmethod
    def total_by_year(
        self, user_name: str, from_year: int, to_year: int, tag_key: Optional[str]
    ) -> List[YearTotal]:
        """Totals grouped by year across the inclusive range, every year present
        (zero-filled), restricted to expenses carrying `tag_key` when given."""
        ...
