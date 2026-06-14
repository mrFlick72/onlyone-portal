from pathlib import Path
from typing import Iterator, List

import psycopg
import pytest
from testcontainers.postgres import PostgresContainer

from app.analytic.adapter.db.repository import PostgresExpenseProjectionRepository
from app.analytic.domain.expense import ExpenseTag, ProjectedExpense
from app.money.domain.money import Money
from app.time.domain.date import Date

pytestmark = pytest.mark.integration

# tests/analytic/adapter/db/ -> analytic-api/scripts/init.sql
INIT_SQL = Path(__file__).parents[4] / "scripts" / "init.sql"


@pytest.fixture(scope="module")
def database_url() -> Iterator[str]:
    with PostgresContainer("postgres:16-alpine") as container:
        url = container.get_connection_url().replace("+psycopg2", "")
        with psycopg.connect(url) as conn, conn.cursor() as cur:
            cur.execute(INIT_SQL.read_text())
        yield url


@pytest.fixture(scope="module")
def repository(database_url: str) -> Iterator[PostgresExpenseProjectionRepository]:
    repo = PostgresExpenseProjectionRepository(database_url)
    repo.open()
    yield repo
    repo.close()


@pytest.fixture(autouse=True)
def clean(database_url: str) -> None:
    with psycopg.connect(database_url) as conn, conn.cursor() as cur:
        cur.execute("TRUNCATE budget_expense_projection, budget_expense_tag")


def expense(
    expense_id: str,
    user: str,
    date: str,
    amount: str,
    tags: List[tuple[str, str]],
    note: str = "",
) -> ProjectedExpense:
    return ProjectedExpense(
        id=expense_id,
        user_name=user,
        date=Date.date_for(date),
        amount=Money.money_for(amount),
        note=note,
        tags=[ExpenseTag(key=key, value=value) for key, value in tags],
    )


def tag_pairs(repo: PostgresExpenseProjectionRepository, **kwargs: object) -> list[tuple[str, str]]:
    return [(t.tag, t.total.stringify_amount()) for t in repo.total_by_tag(**kwargs)]  # type: ignore[arg-type]


def year_pairs(repo: PostgresExpenseProjectionRepository, **kwargs: object) -> list[tuple[int, str]]:
    return [(t.year.content, t.total.stringify_amount()) for t in repo.total_by_year(**kwargs)]  # type: ignore[arg-type]


def test_total_by_tag_groups_by_tag_value(repository: PostgresExpenseProjectionRepository) -> None:
    repository.save(expense("1", "jane", "01/02/2026", "10.10", [("food", "Food")]))
    repository.save(expense("2", "jane", "15/02/2026", "0.20", [("food", "Food")]))
    repository.save(expense("3", "jane", "20/02/2026", "5.00", [("car", "Car")]))

    assert tag_pairs(repository, user_name="jane", year=2026, month=None, tag_keys=[]) == [
        ("Car", "5.00"),
        ("Food", "10.30"),
    ]


def test_total_by_tag_counts_multi_tag_expense_in_every_bucket(
    repository: PostgresExpenseProjectionRepository,
) -> None:
    repository.save(
        expense("1", "jane", "01/02/2026", "7.50", [("food", "Food"), ("car", "Car")])
    )

    assert tag_pairs(repository, user_name="jane", year=2026, month=None, tag_keys=[]) == [
        ("Car", "7.50"),
        ("Food", "7.50"),
    ]


def test_total_by_tag_filters_by_tag_key_but_keeps_sibling_tags(
    repository: PostgresExpenseProjectionRepository,
) -> None:
    # an expense matched by the food filter still contributes all of its tags,
    # preserving the previous REST-pull behaviour
    repository.save(
        expense("1", "jane", "01/02/2026", "7.50", [("food", "Food"), ("car", "Car")])
    )
    repository.save(expense("2", "jane", "02/02/2026", "3.00", [("rent", "Rent")]))

    assert tag_pairs(
        repository, user_name="jane", year=2026, month=None, tag_keys=["food"]
    ) == [
        ("Car", "7.50"),
        ("Food", "7.50"),
    ]


def test_total_by_tag_honours_month_filter(repository: PostgresExpenseProjectionRepository) -> None:
    repository.save(expense("1", "jane", "01/02/2026", "10.00", [("food", "Food")]))
    repository.save(expense("2", "jane", "01/03/2026", "5.00", [("food", "Food")]))

    assert tag_pairs(repository, user_name="jane", year=2026, month=2, tag_keys=[]) == [
        ("Food", "10.00"),
    ]


def test_total_by_tag_is_scoped_per_user(repository: PostgresExpenseProjectionRepository) -> None:
    repository.save(expense("1", "jane", "01/02/2026", "10.00", [("food", "Food")]))
    repository.save(expense("2", "bob", "01/02/2026", "99.00", [("food", "Food")]))

    assert tag_pairs(repository, user_name="jane", year=2026, month=None, tag_keys=[]) == [
        ("Food", "10.00"),
    ]


def test_save_is_idempotent_upsert(repository: PostgresExpenseProjectionRepository) -> None:
    repository.save(expense("1", "jane", "01/02/2026", "10.00", [("food", "Food")]))
    repository.save(expense("1", "jane", "01/02/2026", "25.00", [("food", "Food")]))

    assert tag_pairs(repository, user_name="jane", year=2026, month=None, tag_keys=[]) == [
        ("Food", "25.00"),
    ]


def test_update_replaces_tags(repository: PostgresExpenseProjectionRepository) -> None:
    repository.save(expense("1", "jane", "01/02/2026", "10.00", [("food", "Food")]))
    repository.save(expense("1", "jane", "01/02/2026", "10.00", [("car", "Car")]))

    assert tag_pairs(repository, user_name="jane", year=2026, month=None, tag_keys=[]) == [
        ("Car", "10.00"),
    ]


def test_delete_removes_expense_from_totals(repository: PostgresExpenseProjectionRepository) -> None:
    repository.save(expense("1", "jane", "01/02/2026", "10.00", [("food", "Food")]))
    repository.delete("1")

    assert tag_pairs(repository, user_name="jane", year=2026, month=None, tag_keys=[]) == []


def test_delete_is_a_noop_for_unknown_id(repository: PostgresExpenseProjectionRepository) -> None:
    repository.delete("does-not-exist")


def test_total_by_year_groups_and_zero_fills(repository: PostgresExpenseProjectionRepository) -> None:
    repository.save(expense("1", "jane", "01/02/2024", "10.00", [("food", "Food")]))
    repository.save(expense("2", "jane", "01/03/2024", "2.50", [("car", "Car")]))
    repository.save(expense("3", "jane", "01/02/2026", "1.00", [("food", "Food")]))

    assert year_pairs(repository, user_name="jane", from_year=2024, to_year=2026, tag_key=None) == [
        (2024, "12.50"),
        (2025, "0.00"),
        (2026, "1.00"),
    ]


def test_total_by_year_counts_each_expense_once(repository: PostgresExpenseProjectionRepository) -> None:
    repository.save(
        expense("1", "jane", "01/02/2026", "7.50", [("food", "Food"), ("car", "Car")])
    )

    assert year_pairs(repository, user_name="jane", from_year=2026, to_year=2026, tag_key=None) == [
        (2026, "7.50"),
    ]


def test_total_by_year_filters_by_tag_key(repository: PostgresExpenseProjectionRepository) -> None:
    repository.save(expense("1", "jane", "01/02/2026", "10.00", [("food", "Food")]))
    repository.save(expense("2", "jane", "01/02/2026", "5.00", [("car", "Car")]))

    assert year_pairs(repository, user_name="jane", from_year=2026, to_year=2026, tag_key="food") == [
        (2026, "10.00"),
    ]
