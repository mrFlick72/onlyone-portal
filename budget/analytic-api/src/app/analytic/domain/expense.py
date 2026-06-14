from dataclasses import dataclass, field
from typing import List

from app.money.domain.money import Money
from app.time.domain.date import Date
from app.time.domain.year import Year


@dataclass
class ExpenseTag:
    key: str
    value: str


@dataclass
class ProjectedExpense:
    id: str
    user_name: str
    date: Date
    amount: Money
    note: str
    tags: List[ExpenseTag] = field(default_factory=list)


@dataclass
class TagTotal:
    tag: str
    total: Money


@dataclass
class YearTotal:
    year: Year
    total: Money
