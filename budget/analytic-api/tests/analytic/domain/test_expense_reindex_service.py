from typing import List, Optional

from app.analytic.domain.expense import ProjectedExpense, TagTotal, YearTotal
from app.analytic.domain.reindex import BudgetExpenseSource, ExpenseReindexService
from app.analytic.domain.repository import ExpenseProjectionRepository
from app.money.domain.money import Money
from app.time.domain.date import Date


def expense(expense_id: str) -> ProjectedExpense:
    return ProjectedExpense(
        id=expense_id,
        user_name="jane",
        date=Date.date_for("01/02/2026"),
        amount=Money.money_for("10.00"),
        note="",
        tags=[],
    )


class FakeSource(BudgetExpenseSource):
    def __init__(self, expenses: List[ProjectedExpense]) -> None:
        self.expenses = expenses
        self.calls: List[tuple] = []

    def fetch_expenses(self, from_year: int, to_year: int) -> List[ProjectedExpense]:
        self.calls.append((from_year, to_year))
        return self.expenses


class RecordingRepository(ExpenseProjectionRepository):
    def __init__(self) -> None:
        self.saved: List[ProjectedExpense] = []

    def save(self, expense: ProjectedExpense) -> None:
        self.saved.append(expense)

    def delete(self, expense_id: str) -> None: ...

    def total_by_tag(
        self, user_name: str, year: int, month: Optional[int], tag_keys: List[str]
    ) -> List[TagTotal]:
        return []

    def total_by_year(
        self, user_name: str, from_year: int, to_year: int, tag_key: Optional[str]
    ) -> List[YearTotal]:
        return []


def test_reindex_saves_every_fetched_expense_and_returns_count() -> None:
    source = FakeSource([expense("1"), expense("2"), expense("3")])
    repository = RecordingRepository()

    imported = ExpenseReindexService(source, repository).reindex(2024, 2026)

    assert imported == 3
    assert [e.id for e in repository.saved] == ["1", "2", "3"]
    assert source.calls == [(2024, 2026)]


def test_reindex_with_no_expenses_returns_zero() -> None:
    source = FakeSource([])
    repository = RecordingRepository()

    imported = ExpenseReindexService(source, repository).reindex(2026, 2026)

    assert imported == 0
    assert repository.saved == []
