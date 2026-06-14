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


class BudgetExpenseEventHandler:
    """Applies a decoded `budget-api.expense` event to the projection. Kept free
    of any Kafka type so it can be unit-tested with a plain dict."""

    def __init__(self, repository: ExpenseProjectionRepository) -> None:
        self.repository = repository

    def handle(self, event: Dict[str, Any]) -> None:
        action = event["action"]
        payload = event["payload"]

        if action in (CREATE_ACTION, UPDATE_ACTION):
            self.repository.save(self._to_expense(payload))
        elif action == DELETE_ACTION:
            self.repository.delete(payload["id"])
        else:
            logger.warning("Ignoring budget-expense event with unknown action: %s", action)

    @staticmethod
    def _to_expense(payload: Dict[str, Any]) -> ProjectedExpense:
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
            try:
                self.handler.handle(json.loads(message.value()))
                self._consumer.commit(message=message, asynchronous=False)
            except Exception:
                # leave the offset uncommitted so the message is redelivered
                logger.exception("Failed to apply budget-expense event; will retry")
