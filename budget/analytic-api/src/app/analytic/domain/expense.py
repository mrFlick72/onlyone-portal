from dataclasses import dataclass, field
from typing import List

from app.money.domain.money import Money
from app.time.domain.date import Date
from app.time.domain.year import Year


@dataclass
class ExpenseRecord:
    id: str
    date: Date
    amount: Money
    note: str
    tag_values: List[str] = field(default_factory=list)


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
