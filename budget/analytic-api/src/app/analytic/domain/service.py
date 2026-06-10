from dataclasses import dataclass
from typing import Dict, List, Optional
from app.analytic.domain.expense import ExpenseRecord
from app.money.domain.money import Money
from app.time.domain.year import Year
from abc import ABC, abstractmethod

ALL_MONTHS = range(1, 13)


@dataclass
class BudgetExpenseAnalysisRequest:
    year: int
    month: int
    tags: List[str]


@dataclass
class TagTotal:
    tag: str
    total: Money


@dataclass
class YearTotal:
    year: Year
    total: Money


class ExpenseLoader(ABC):

    def __init__(self) -> None:
        super(self)

    @abstractmethod
    def expenseFor(self, request: BudgetExpenseAnalysisRequest) -> List[ExpenseRecord]: ...

    @abstractmethod
    def expenses_for_all(
        self, requests: List[BudgetExpenseAnalysisRequest]
    ) -> List[ExpenseRecord]: ...


class BudgetExpenseAnalysisService:

    def __init__(self, expense_loader: ExpenseLoader) -> None:
        self.expense_loader = expense_loader

    def total_by_tag(
        self, year: int, month: Optional[int], tags: List[str]
    ) -> List[TagTotal]:
        months = [month] if month is not None else list(ALL_MONTHS)
        records = self.expense_loader.expenses_for_all(
            [BudgetExpenseAnalysisRequest(year, m, tags) for m in months]
        )

        totals: Dict[str, Money] = {}
        for record in records:
            amount = _scaled_money_for(record.amount)
            for tag in record.tag_values:
                accumulated = totals.get(tag)
                if accumulated is None:
                    totals[tag] = amount
                else:
                    totals[tag] = Money(accumulated.plus(amount))

        return [TagTotal(tag=tag, total=totals[tag]) for tag in sorted(totals)]

    def total_by_year(
        self, from_year: int, to_year: int, tag: Optional[str]
    ) -> List[YearTotal]:
        years = range(from_year, to_year + 1)
        records = self.expense_loader.expenses_for_all(
            [
                BudgetExpenseAnalysisRequest(y, m, [tag] if tag is not None else [])
                for y in years
                for m in ALL_MONTHS
            ]
        )

        totals: Dict[int, Money] = {y: Money.money_for("0.00") for y in years}
        for record in records:
            year = record.date.content.year
            totals[year] = Money(totals[year].plus(_scaled_money_for(record.amount)))

        return [YearTotal(year=Year(y), total=totals[y]) for y in years]


def _scaled_money_for(amount: float) -> Money:
    return Money.money_for(f"{amount:.2f}")
