from typing import Generator, List, Optional

import pytest
from dependency_injector import providers
from fastapi.testclient import TestClient

from app.analytic.domain.service import TagTotal, YearTotal
from app.money.domain.money import Money
from app.server import app, application_container
from app.time.domain.year import Year


class FakeBudgetExpenseAnalysisService:

    def __init__(
        self,
        tag_totals: Optional[List[TagTotal]] = None,
        year_totals: Optional[List[YearTotal]] = None,
    ) -> None:
        self.tag_totals = tag_totals or []
        self.year_totals = year_totals or []
        self.total_by_tag_calls: List[tuple] = []
        self.total_by_year_calls: List[tuple] = []

    def total_by_tag(
        self, year: int, month: Optional[int], tags: List[str]
    ) -> List[TagTotal]:
        self.total_by_tag_calls.append((year, month, tags))
        return self.tag_totals

    def total_by_year(
        self, from_year: int, to_year: int, tag: Optional[str]
    ) -> List[YearTotal]:
        self.total_by_year_calls.append((from_year, to_year, tag))
        return self.year_totals


@pytest.fixture
def client() -> TestClient:
    return TestClient(app)


@pytest.fixture
def fake_service() -> Generator[FakeBudgetExpenseAnalysisService, None, None]:
    service = FakeBudgetExpenseAnalysisService(
        tag_totals=[
            TagTotal(tag="car", total=Money.money_for("5.00")),
            TagTotal(tag="food", total=Money.money_for("10.30")),
        ],
        year_totals=[
            YearTotal(year=Year(2025), total=Money.money_for("0.00")),
            YearTotal(year=Year(2026), total=Money.money_for("100.50")),
        ],
    )
    provider = application_container.analytic_config_container.budget_expense_analysis_service
    provider.override(providers.Object(service))
    yield service
    provider.reset_override()


def test_total_by_tag_returns_totals(client, fake_service):
    response = client.put(
        "/api/analytic/budget/expense/total-by-tag",
        json={"year": 2026, "month": 2, "tags": ["food", "car"]},
    )

    assert response.status_code == 200
    assert response.json() == [
        {"tag": "car", "total": "5.00"},
        {"tag": "food", "total": "10.30"},
    ]
    assert fake_service.total_by_tag_calls == [(2026, 2, ["food", "car"])]


def test_total_by_tag_month_and_tags_are_optional(client, fake_service):
    response = client.put(
        "/api/analytic/budget/expense/total-by-tag", json={"year": 2026}
    )

    assert response.status_code == 200
    assert fake_service.total_by_tag_calls == [(2026, None, [])]


def test_total_by_tag_rejects_invalid_month(client, fake_service):
    response = client.put(
        "/api/analytic/budget/expense/total-by-tag",
        json={"year": 2026, "month": 13, "tags": []},
    )

    assert response.status_code == 422


def test_total_by_year_returns_totals(client, fake_service):
    response = client.put(
        "/api/analytic/budget/expense/total-by-year",
        json={"fromYear": 2025, "toYear": 2026, "tag": "food"},
    )

    assert response.status_code == 200
    assert response.json() == [
        {"year": 2025, "total": "0.00"},
        {"year": 2026, "total": "100.50"},
    ]
    assert fake_service.total_by_year_calls == [(2025, 2026, "food")]


def test_total_by_year_tag_is_optional(client, fake_service):
    response = client.put(
        "/api/analytic/budget/expense/total-by-year",
        json={"fromYear": 2022, "toYear": 2026},
    )

    assert response.status_code == 200
    assert fake_service.total_by_year_calls == [(2022, 2026, None)]


def test_total_by_year_rejects_inverted_range(client, fake_service):
    response = client.put(
        "/api/analytic/budget/expense/total-by-year",
        json={"fromYear": 2026, "toYear": 2022},
    )

    assert response.status_code == 422


def test_total_by_year_rejects_too_wide_range(client, fake_service):
    response = client.put(
        "/api/analytic/budget/expense/total-by-year",
        json={"fromYear": 2000, "toYear": 2026},
    )

    assert response.status_code == 422
