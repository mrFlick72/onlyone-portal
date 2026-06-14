from typing import List, Optional

from app.analytic.domain.expense import TagTotal, YearTotal
from app.analytic.domain.repository import ExpenseProjectionRepository
from app.infrastructure.security.security_context_resolver import (
    SecurityContextResolver,
)


class BudgetExpenseAnalysisService:
    """Answers the analytics queries from the local projection, scoping every
    read to the user resolved from the current security context."""

    def __init__(
        self,
        repository: ExpenseProjectionRepository,
        security_context_resolver: SecurityContextResolver,
    ) -> None:
        self.repository = repository
        self.security_context_resolver = security_context_resolver

    def total_by_tag(
        self, year: int, month: Optional[int], tags: List[str]
    ) -> List[TagTotal]:
        return self.repository.total_by_tag(
            self._current_user_name(), year, month, tags
        )

    def total_by_year(
        self, from_year: int, to_year: int, tag: Optional[str]
    ) -> List[YearTotal]:
        return self.repository.total_by_year(
            self._current_user_name(), from_year, to_year, tag
        )

    def _current_user_name(self) -> str:
        return self.security_context_resolver.get_security_context().user_name.content
