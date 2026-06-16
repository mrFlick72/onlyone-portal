import json
from typing import List, Optional
from unittest.mock import MagicMock, patch

import pytest

from app.analytic.adapter.kafka.consumer import (
    BudgetExpenseConsumer,
    BudgetExpenseEventHandler,
    MalformedEventError,
)
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


def test_missing_tags_yields_no_tags() -> None:
    repository = RecordingRepository()
    payload = event("CREATE")
    del payload["payload"]["tags"]

    BudgetExpenseEventHandler(repository).handle(payload)

    assert repository.saved[0].tags == []


@pytest.mark.parametrize(
    "bad_event",
    [
        {"action": "ARCHIVE", "payload": {"id": "exp-1"}},  # unknown action
        {"payload": {"id": "exp-1"}},                        # missing action
        {"action": "CREATE"},                                # missing payload
        {"action": "DELETE", "payload": {}},                 # delete without id
        {"action": "CREATE", "payload": {"id": "exp-1"}},    # incomplete payload (no date/amount)
    ],
)
def test_malformed_events_raise_and_touch_nothing(bad_event: dict) -> None:
    repository = RecordingRepository()

    with pytest.raises(MalformedEventError):
        BudgetExpenseEventHandler(repository).handle(bad_event)

    assert repository.saved == []
    assert repository.deleted == []


class _Message:
    def __init__(self, value: bytes) -> None:
        self._value = value

    def value(self) -> bytes:
        return self._value


def _consumer_with(repository: ExpenseProjectionRepository) -> tuple[BudgetExpenseConsumer, MagicMock]:
    with patch("app.analytic.adapter.kafka.consumer.Consumer") as mock_consumer_cls:
        consumer = BudgetExpenseConsumer(
            BudgetExpenseEventHandler(repository), "broker:9092", "topic", "group"
        )
    return consumer, mock_consumer_cls.return_value


def test_consume_applies_and_commits_valid_event() -> None:
    repository = RecordingRepository()
    consumer, kafka = _consumer_with(repository)

    consumer._consume(_Message(json.dumps(event("CREATE")).encode()))

    assert len(repository.saved) == 1
    kafka.commit.assert_called_once()


def test_consume_skips_and_commits_invalid_json() -> None:
    repository = RecordingRepository()
    consumer, kafka = _consumer_with(repository)

    consumer._consume(_Message(b"not-json"))

    assert repository.saved == []
    kafka.commit.assert_called_once()


def test_consume_skips_and_commits_malformed_event() -> None:
    repository = RecordingRepository()
    consumer, kafka = _consumer_with(repository)

    consumer._consume(_Message(json.dumps({"action": "CREATE"}).encode()))

    assert repository.saved == []
    kafka.commit.assert_called_once()


def test_consume_does_not_commit_on_transient_repository_failure() -> None:
    class FailingRepository(RecordingRepository):
        def save(self, expense: ProjectedExpense) -> None:
            raise RuntimeError("database is down")

    consumer, kafka = _consumer_with(FailingRepository())

    consumer._consume(_Message(json.dumps(event("CREATE")).encode()))

    kafka.commit.assert_not_called()
