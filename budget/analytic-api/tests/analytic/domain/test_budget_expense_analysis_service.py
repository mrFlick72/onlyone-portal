from typing import List, Optional, Tuple

from app.analytic.domain.expense import ProjectedExpense, TagTotal, YearTotal
from app.analytic.domain.repository import ExpenseProjectionRepository
from app.analytic.domain.service import BudgetExpenseAnalysisService
from app.infrastructure.security.model import SecurityContext, UserName
from app.infrastructure.security.security_context_resolver import (
    SecurityContextResolver,
)
from app.money.domain.money import Money
from app.time.domain.year import Year


class FakeProjectionRepository(ExpenseProjectionRepository):
    def __init__(self) -> None:
        self.tag_calls: List[Tuple[str, int, Optional[int], List[str]]] = []
        self.year_calls: List[Tuple[str, int, int, Optional[str]]] = []
        self.tag_totals = [TagTotal(tag="food", total=Money.money_for("10.30"))]
        self.year_totals = [YearTotal(year=Year(2026), total=Money.money_for("1.00"))]

    def save(self, expense: ProjectedExpense) -> None: ...

    def delete(self, expense_id: str) -> None: ...

    def total_by_tag(
        self, user_name: str, year: int, month: Optional[int], tag_keys: List[str]
    ) -> List[TagTotal]:
        self.tag_calls.append((user_name, year, month, tag_keys))
        return self.tag_totals

    def total_by_year(
        self, user_name: str, from_year: int, to_year: int, tag_key: Optional[str]
    ) -> List[YearTotal]:
        self.year_calls.append((user_name, from_year, to_year, tag_key))
        return self.year_totals


class FakeSecurityContextResolver(SecurityContextResolver):
    def __init__(self, user_name: str) -> None:
        self.context = SecurityContext(token="a-token", user_name=UserName(user_name))

    def get_security_context(self) -> SecurityContext:
        return self.context

    def set_security_context(self, ctx: SecurityContext) -> None:
        self.context = ctx


def service_for(
    user: str = "jane",
) -> Tuple[BudgetExpenseAnalysisService, FakeProjectionRepository]:
    repository = FakeProjectionRepository()
    service = BudgetExpenseAnalysisService(repository, FakeSecurityContextResolver(user))
    return service, repository


def test_total_by_tag_scopes_to_current_user_and_delegates() -> None:
    service, repository = service_for("jane")

    result = service.total_by_tag(2026, 2, ["food", "car"])

    assert result == repository.tag_totals
    assert repository.tag_calls == [("jane", 2026, 2, ["food", "car"])]


def test_total_by_tag_passes_optional_month_through() -> None:
    service, repository = service_for("jane")

    service.total_by_tag(2026, None, [])

    assert repository.tag_calls == [("jane", 2026, None, [])]


def test_total_by_year_scopes_to_current_user_and_delegates() -> None:
    service, repository = service_for("bob")

    result = service.total_by_year(2024, 2026, "food")

    assert result == repository.year_totals
    assert repository.year_calls == [("bob", 2024, 2026, "food")]
