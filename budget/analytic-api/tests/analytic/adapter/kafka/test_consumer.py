from typing import List, Optional

from app.analytic.adapter.kafka.consumer import BudgetExpenseEventHandler
from app.analytic.domain.expense import ProjectedExpense, TagTotal, YearTotal
from app.analytic.domain.repository import ExpenseProjectionRepository


class RecordingRepository(ExpenseProjectionRepository):
    def __init__(self) -> None:
        self.saved: List[ProjectedExpense] = []
        self.deleted: List[str] = []

    def save(self, expense: ProjectedExpense) -> None:
        self.saved.append(expense)

    def delete(self, expense_id: str) -> None:
        self.deleted.append(expense_id)

    def total_by_tag(
        self, user_name: str, year: int, month: Optional[int], tag_keys: List[str]
    ) -> List[TagTotal]:
        return []

    def total_by_year(
        self, user_name: str, from_year: int, to_year: int, tag_key: Optional[str]
    ) -> List[YearTotal]:
        return []


def event(action: str) -> dict:
    return {
        "action": action,
        "payload": {
            "id": "exp-1",
            "userName": "jane",
            "date": "15/02/2026",
            "amount": "12.50",
            "note": "groceries",
            "tags": [
                {"key": "food", "value": "Food"},
                {"key": "car", "value": "Car"},
            ],
        },
    }


def test_create_event_saves_projected_expense() -> None:
    repository = RecordingRepository()
    BudgetExpenseEventHandler(repository).handle(event("CREATE"))

    assert repository.deleted == []
    assert len(repository.saved) == 1
    saved = repository.saved[0]
    assert saved.id == "exp-1"
    assert saved.user_name == "jane"
    assert saved.date.formatted_date() == "15/02/2026"
    assert saved.amount.stringify_amount() == "12.50"
    assert saved.note == "groceries"
    assert [(tag.key, tag.value) for tag in saved.tags] == [
        ("food", "Food"),
        ("car", "Car"),
    ]


def test_update_event_saves_projected_expense() -> None:
    repository = RecordingRepository()
    BudgetExpenseEventHandler(repository).handle(event("UPDATE"))

    assert len(repository.saved) == 1
    assert repository.deleted == []


def test_delete_event_removes_by_id() -> None:
    repository = RecordingRepository()
    BudgetExpenseEventHandler(repository).handle(event("DELETE"))

    assert repository.saved == []
    assert repository.deleted == ["exp-1"]


def test_unknown_action_is_ignored() -> None:
    repository = RecordingRepository()
    BudgetExpenseEventHandler(repository).handle(event("ARCHIVE"))

    assert repository.saved == []
    assert repository.deleted == []


def test_missing_tags_yields_no_tags() -> None:
    repository = RecordingRepository()
    payload = event("CREATE")
    del payload["payload"]["tags"]

    BudgetExpenseEventHandler(repository).handle(payload)

    assert repository.saved[0].tags == []
