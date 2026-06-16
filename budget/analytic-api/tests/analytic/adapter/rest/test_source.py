import json
from typing import List

import httpx

from app.analytic.adapter.rest.source import RestBudgetExpenseSource
from app.infrastructure.security.model import SecurityContext, UserName
from app.infrastructure.security.security_context_resolver import (
    SecurityContextResolver,
)

BASE_URL = "http://budget-api.test"


class FixedResolver(SecurityContextResolver):
    def __init__(self, token: str, user_name: str) -> None:
        self.context = SecurityContext(token=token, user_name=UserName(user_name))

    def get_security_context(self) -> SecurityContext:
        return self.context

    def set_security_context(self, ctx: SecurityContext) -> None:
        self.context = ctx


def expense_payload() -> dict:
    return {
        "dailyBudgetExpenseRepresentationList": [
            {
                "date": "15/02/2026",
                "budgetExpenseRepresentationList": [
                    {
                        "id": "exp-1",
                        "date": "15/02/2026",
                        "amount": "12.50",
                        "note": "groceries",
                        "tags": [
                            {"tagKey": "food", "tagValue": "Food"},
                            {"tagKey": "car", "tagValue": "Car"},
                        ],
                    }
                ],
            }
        ]
    }


def source_with_handler(handler) -> RestBudgetExpenseSource:
    return RestBudgetExpenseSource(
        budget_api_base_url=BASE_URL,
        security_context_resolver=FixedResolver("a-token", "jane"),
        transport=httpx.MockTransport(handler),
    )


def test_fetch_expenses_queries_every_month_in_range_and_propagates_token() -> None:
    seen: List[tuple] = []

    def handler(request: httpx.Request) -> httpx.Response:
        body = json.loads(request.content)
        seen.append((body["year"], body["month"], request.headers.get("authorization")))
        return httpx.Response(200, json={"dailyBudgetExpenseRepresentationList": []})

    source_with_handler(handler).fetch_expenses(2025, 2026)

    assert len(seen) == 24  # two years x twelve months
    assert {token for _, _, token in seen} == {"Bearer a-token"}
    assert ("2026", "2", "Bearer a-token") in seen


def test_fetch_expenses_maps_response_to_projected_expenses_with_tag_key_and_value() -> None:
    def handler(request: httpx.Request) -> httpx.Response:
        body = json.loads(request.content)
        if (body["year"], body["month"]) == ("2026", "2"):
            return httpx.Response(200, json=expense_payload())
        return httpx.Response(200, json={"dailyBudgetExpenseRepresentationList": []})

    expenses = source_with_handler(handler).fetch_expenses(2026, 2026)

    assert len(expenses) == 1
    projected = expenses[0]
    assert projected.id == "exp-1"
    assert projected.user_name == "jane"
    assert projected.date.formatted_date() == "15/02/2026"
    assert projected.amount.stringify_amount() == "12.50"
    assert projected.note == "groceries"
    assert [(tag.key, tag.value) for tag in projected.tags] == [
        ("food", "Food"),
        ("car", "Car"),
    ]
