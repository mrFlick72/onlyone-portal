from decimal import Decimal
from typing import Dict, List, Optional

from psycopg_pool import ConnectionPool

from app.analytic.domain.expense import ProjectedExpense, TagTotal, YearTotal
from app.analytic.domain.repository import ExpenseProjectionRepository
from app.money.domain.money import Money
from app.time.domain.year import Year


class PostgresExpenseProjectionRepository(ExpenseProjectionRepository):
    """Postgres-backed projection. Owns its own connection pool, opened lazily on
    `open()` (from the FastAPI lifespan) and never during tests.

    The schema (`scripts/init.sql`) is applied out-of-band — by the local
    docker-compose Postgres init, by deployment migrations, or by tests — not by
    this adapter."""

    def __init__(self, database_url: str) -> None:
        self._pool = ConnectionPool(database_url, min_size=1, max_size=10, open=False)

    def open(self) -> None:
        self._pool.open(wait=True)

    def close(self) -> None:
        self._pool.close()

    def save(self, expense: ProjectedExpense) -> None:
        # de-duplicate by key: the tag table is keyed on (expense_id, tag_key)
        tags = list({tag.key: tag for tag in expense.tags}.values())
        with self._pool.connection() as conn, conn.cursor() as cur:
            cur.execute(
                """
                INSERT INTO budget_expense_projection (id, user_name, expense_date, amount, note)
                VALUES (%s, %s, %s, %s, %s)
                ON CONFLICT (id) DO UPDATE SET
                    user_name = EXCLUDED.user_name,
                    expense_date = EXCLUDED.expense_date,
                    amount = EXCLUDED.amount,
                    note = EXCLUDED.note
                """,
                (
                    expense.id,
                    expense.user_name,
                    expense.date.content,
                    expense.amount.amount,
                    expense.note,
                ),
            )
            cur.execute(
                "DELETE FROM budget_expense_tag WHERE expense_id = %s", (expense.id,)
            )
            if tags:
                cur.executemany(
                    "INSERT INTO budget_expense_tag (expense_id, tag_key, tag_value) VALUES (%s, %s, %s)",
                    [(expense.id, tag.key, tag.value) for tag in tags],
                )

    def delete(self, expense_id: str) -> None:
        with self._pool.connection() as conn, conn.cursor() as cur:
            cur.execute(
                "DELETE FROM budget_expense_projection WHERE id = %s", (expense_id,)
            )

    def total_by_tag(
        self, user_name: str, year: int, month: Optional[int], tag_keys: List[str]
    ) -> List[TagTotal]:
        clauses = [
            "e.user_name = %(user)s",
            "EXTRACT(YEAR FROM e.expense_date) = %(year)s",
        ]
        params: Dict[str, object] = {"user": user_name, "year": year}
        if month is not None:
            clauses.append("EXTRACT(MONTH FROM e.expense_date) = %(month)s")
            params["month"] = month
        if tag_keys:
            clauses.append(
                "EXISTS (SELECT 1 FROM budget_expense_tag f "
                "WHERE f.expense_id = e.id AND f.tag_key = ANY(%(tag_keys)s))"
            )
            params["tag_keys"] = tag_keys

        query = (
            "SELECT t.tag_value, SUM(e.amount) "
            "FROM budget_expense_projection e "
            "JOIN budget_expense_tag t ON t.expense_id = e.id "
            f"WHERE {' AND '.join(clauses)} "
            "GROUP BY t.tag_value ORDER BY t.tag_value"
        )

        with self._pool.connection() as conn, conn.cursor() as cur:
            cur.execute(query, params)
            rows = cur.fetchall()

        return [
            TagTotal(tag=str(tag_value), total=self._money(total))
            for tag_value, total in rows
        ]

    def total_by_year(
        self, user_name: str, from_year: int, to_year: int, tag_key: Optional[str]
    ) -> List[YearTotal]:
        clauses = [
            "e.user_name = %(user)s",
            "EXTRACT(YEAR FROM e.expense_date) BETWEEN %(from_year)s AND %(to_year)s",
        ]
        params: Dict[str, object] = {
            "user": user_name,
            "from_year": from_year,
            "to_year": to_year,
        }
        if tag_key is not None:
            clauses.append(
                "EXISTS (SELECT 1 FROM budget_expense_tag f "
                "WHERE f.expense_id = e.id AND f.tag_key = %(tag_key)s)"
            )
            params["tag_key"] = tag_key

        query = (
            "SELECT EXTRACT(YEAR FROM e.expense_date)::int, SUM(e.amount) "
            "FROM budget_expense_projection e "
            f"WHERE {' AND '.join(clauses)} "
            "GROUP BY 1"
        )

        with self._pool.connection() as conn, conn.cursor() as cur:
            cur.execute(query, params)
            totals_by_year = {int(year): self._money(total) for year, total in cur.fetchall()}

        return [
            YearTotal(
                year=Year(year),
                total=totals_by_year.get(year, Money.money_for("0.00")),
            )
            for year in range(from_year, to_year + 1)
        ]

    @staticmethod
    def _money(amount: Decimal) -> Money:
        return Money.money_for(str(amount))
