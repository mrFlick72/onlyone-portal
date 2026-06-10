from decimal import Decimal
from typing import List

from app.analytic.domain.expense import ExpenseRecord
from app.analytic.domain.service import (
    BudgetExpenseAnalysisRequest,
    BudgetExpenseAnalysisService,
    ExpenseLoader,
)
from app.time.domain.date import Date


class FakeExpenseLoader(ExpenseLoader):

    def __init__(self, records: List[ExpenseRecord]) -> None:
        self.records = records
        self.received_requests: List[BudgetExpenseAnalysisRequest] = []

    def expenseFor(
        self, request: BudgetExpenseAnalysisRequest
    ) -> List[ExpenseRecord]:
        self.received_requests.append(request)
        return self.records

    def expenses_for_all(
        self, requests: List[BudgetExpenseAnalysisRequest]
    ) -> List[ExpenseRecord]:
        self.received_requests.extend(requests)
        return self.records


def expense_record(date: str, amount: float, tags: List[str]) -> ExpenseRecord:
    return ExpenseRecord(
        id="an-id", date=Date.date_for(date), amount=amount, note="", tag_values=tags
    )


def test_total_by_tag_groups_amounts_by_tag() -> None:
    loader = FakeExpenseLoader(
        [
            expense_record("01/02/2026", 10.10, ["food"]),
            expense_record("15/02/2026", 0.20, ["food"]),
            expense_record("20/02/2026", 5.00, ["car"]),
        ]
    )
    service = BudgetExpenseAnalysisService(loader)

    totals = service.total_by_tag(2026, 2, [])

    assert [(t.tag, t.total.stringify_amount()) for t in totals] == [
        ("car", "5.00"),
        ("food", "10.30"),
    ]


def test_total_by_tag_counts_multi_tag_record_in_every_bucket() -> None:
    loader = FakeExpenseLoader([expense_record("01/02/2026", 7.50, ["food", "car"])])
    service = BudgetExpenseAnalysisService(loader)

    totals = service.total_by_tag(2026, 2, [])

    assert [(t.tag, t.total.stringify_amount()) for t in totals] == [
        ("car", "7.50"),
        ("food", "7.50"),
    ]


def test_total_by_tag_skips_untagged_records() -> None:
    loader = FakeExpenseLoader([expense_record("01/02/2026", 7.50, [])])
    service = BudgetExpenseAnalysisService(loader)

    assert service.total_by_tag(2026, 2, []) == []


def test_total_by_tag_with_month_issues_a_single_request() -> None:
    loader = FakeExpenseLoader([])
    service = BudgetExpenseAnalysisService(loader)

    service.total_by_tag(2026, 2, ["food"])

    assert loader.received_requests == [BudgetExpenseAnalysisRequest(2026, 2, ["food"])]


def test_total_by_tag_without_month_issues_a_request_per_month() -> None:
    loader = FakeExpenseLoader([])
    service = BudgetExpenseAnalysisService(loader)

    service.total_by_tag(2026, None, ["food"])

    assert loader.received_requests == [
        BudgetExpenseAnalysisRequest(2026, month, ["food"]) for month in range(1, 13)
    ]


def test_total_by_tag_sums_with_decimal_precision() -> None:
    loader = FakeExpenseLoader(
        [expense_record("01/02/2026", 0.1, ["food"]) for _ in range(3)]
    )
    service = BudgetExpenseAnalysisService(loader)

    totals = service.total_by_tag(2026, 2, [])

    assert totals[0].total.amount == Decimal("0.30")


def test_total_by_year_groups_amounts_by_year_and_zero_fills_empty_years() -> None:
    loader = FakeExpenseLoader(
        [
            expense_record("01/02/2024", 10.00, ["food"]),
            expense_record("01/03/2024", 2.50, ["car"]),
            expense_record("01/02/2026", 1.00, ["food"]),
        ]
    )
    service = BudgetExpenseAnalysisService(loader)

    totals = service.total_by_year(2024, 2026, None)

    assert [(t.year.content, t.total.stringify_amount()) for t in totals] == [
        (2024, "12.50"),
        (2025, "0.00"),
        (2026, "1.00"),
    ]


def test_total_by_year_issues_a_request_per_month_per_year_with_tag_filter() -> None:
    loader = FakeExpenseLoader([])
    service = BudgetExpenseAnalysisService(loader)

    service.total_by_year(2025, 2026, "food")

    assert loader.received_requests == [
        BudgetExpenseAnalysisRequest(year, month, ["food"])
        for year in (2025, 2026)
        for month in range(1, 13)
    ]


def test_total_by_year_without_tag_requests_all_tags() -> None:
    loader = FakeExpenseLoader([])
    service = BudgetExpenseAnalysisService(loader)

    service.total_by_year(2026, 2026, None)

    assert all(request.tags == [] for request in loader.received_requests)
