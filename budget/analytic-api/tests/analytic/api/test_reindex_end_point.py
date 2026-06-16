from typing import Generator, List, Tuple

import pytest
from dependency_injector import providers
from fastapi.testclient import TestClient

from app.server import app, application_container


class FakeReindexService:
    def __init__(self, imported: int) -> None:
        self.imported = imported
        self.calls: List[Tuple[int, int]] = []

    def reindex(self, from_year: int, to_year: int) -> int:
        self.calls.append((from_year, to_year))
        return self.imported


@pytest.fixture
def client() -> TestClient:
    return TestClient(app)


@pytest.fixture
def fake_service() -> Generator[FakeReindexService, None, None]:
    service = FakeReindexService(imported=7)
    provider = application_container.analytic_config_container.expense_reindex_service
    provider.override(providers.Object(service))
    yield service
    provider.reset_override()


def test_reindex_returns_imported_count(client: TestClient, fake_service: FakeReindexService) -> None:
    response = client.post(
        "/api/analytic/budget/expense/reindex",
        json={"fromYear": 2024, "toYear": 2026},
    )

    assert response.status_code == 200
    assert response.json() == {"imported": 7}
    assert fake_service.calls == [(2024, 2026)]


def test_reindex_rejects_inverted_range(client: TestClient, fake_service: FakeReindexService) -> None:
    response = client.post(
        "/api/analytic/budget/expense/reindex",
        json={"fromYear": 2026, "toYear": 2024},
    )

    assert response.status_code == 422
    assert fake_service.calls == []


def test_reindex_rejects_too_wide_range(client: TestClient, fake_service: FakeReindexService) -> None:
    response = client.post(
        "/api/analytic/budget/expense/reindex",
        json={"fromYear": 2000, "toYear": 2026},
    )

    assert response.status_code == 422
