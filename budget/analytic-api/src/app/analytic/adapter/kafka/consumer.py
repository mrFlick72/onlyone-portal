import json
import logging
import threading
from typing import Any, Dict

from confluent_kafka import Consumer

from app.analytic.domain.expense import ExpenseTag, ProjectedExpense
from app.analytic.domain.repository import ExpenseProjectionRepository
from app.money.domain.money import Money
from app.time.domain.date import Date

logger = logging.getLogger(__name__)

CREATE_ACTION = "CREATE"
UPDATE_ACTION = "UPDATE"
DELETE_ACTION = "DELETE"

POLL_TIMEOUT_SECONDS = 1.0


class MalformedEventError(Exception):
    """A budget-expense event that can never be applied (bad structure, unknown
    action, unparseable payload). Such a message is skipped, not retried."""


class BudgetExpenseEventHandler:
    """Applies a decoded `budget-api.expense` event to the projection. Kept free
    of any Kafka type so it can be unit-tested with a plain dict.

    Raises `MalformedEventError` for permanently-bad input (the caller skips it);
    repository failures propagate unchanged so the caller can retry them."""

    def __init__(self, repository: ExpenseProjectionRepository) -> None:
        self.repository = repository

    def handle(self, event: Dict[str, Any]) -> None:
        if not isinstance(event, dict):
            raise MalformedEventError(f"event is not an object: {type(event).__name__}")

        action = event.get("action")
        payload = event.get("payload")
        if not isinstance(payload, dict):
            raise MalformedEventError(f"missing or invalid payload (action={action!r})")

        if action in (CREATE_ACTION, UPDATE_ACTION):
            self.repository.save(self._to_expense(payload))
        elif action == DELETE_ACTION:
            expense_id = payload.get("id")
            if not expense_id:
                raise MalformedEventError("DELETE event without id")
            self.repository.delete(expense_id)
        else:
            raise MalformedEventError(f"unknown action: {action!r}")

    def _to_expense(self, payload: Dict[str, Any]) -> ProjectedExpense:
        # the payload is a pure data transformation: any failure here is a bad
        # message, not a transient fault, so surface it as MalformedEventError
        try:
            return self._build_expense(payload)
        except Exception as error:
            raise MalformedEventError(f"invalid CREATE/UPDATE payload: {error}") from error

    @staticmethod
    def _build_expense(payload: Dict[str, Any]) -> ProjectedExpense:
        return ProjectedExpense(
            id=payload["id"],
            user_name=payload["userName"],
            date=Date.date_for(payload["date"]),
            amount=Money.money_for(payload["amount"]),
            note=payload.get("note") or "",
            tags=[
                ExpenseTag(key=tag["key"], value=tag["value"])
                for tag in (payload.get("tags") or [])
            ],
        )


class BudgetExpenseConsumer:
    """In-process Kafka consumer. Polls on a background thread and commits offsets
    only after the event has been applied (at-least-once; the projection ops are
    idempotent)."""

    def __init__(
        self,
        handler: BudgetExpenseEventHandler,
        bootstrap_servers: str,
        topic: str,
        group_id: str,
    ) -> None:
        self.handler = handler
        self.topic = topic
        self._consumer = Consumer(
            {
                "bootstrap.servers": bootstrap_servers,
                "group.id": group_id,
                "auto.offset.reset": "earliest",
                "enable.auto.commit": False,
            }
        )
        self._thread = threading.Thread(target=self._run, name="budget-expense-consumer", daemon=True)
        self._stop = threading.Event()

    def start(self) -> None:
        self._consumer.subscribe([self.topic])
        self._thread.start()
        logger.info("Budget-expense Kafka consumer started on topic '%s'", self.topic)

    def stop(self) -> None:
        self._stop.set()
        self._thread.join(timeout=10.0)
        self._consumer.close()
        logger.info("Budget-expense Kafka consumer stopped")

    def _run(self) -> None:
        while not self._stop.is_set():
            message = self._consumer.poll(POLL_TIMEOUT_SECONDS)
            if message is None:
                continue
            if message.error() is not None:
                logger.error("Kafka consume error: %s", message.error())
                continue
            self._consume(message)

    def _consume(self, message: Any) -> None:
        try:
            event = json.loads(message.value())
            self.handler.handle(event)
        except (MalformedEventError, json.JSONDecodeError) as error:
            # permanently unprocessable: commit so the poison message cannot block
            # the partition, but never apply it
            logger.warning("Skipping unprocessable budget-expense event: %s", error)
            self._consumer.commit(message=message, asynchronous=False)
        except Exception:
            # transient failure (e.g. database): leave uncommitted to redeliver
            logger.exception("Failed to apply budget-expense event; will retry")
        else:
            self._consumer.commit(message=message, asynchronous=False)
