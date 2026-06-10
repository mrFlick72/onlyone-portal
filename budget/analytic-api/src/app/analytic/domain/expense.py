from dataclasses import dataclass, field
from typing import List

from app.time.domain.date import Date


@dataclass
class ExpenseRecord:
    id: str
    date: Date
    amount: float
    note: str
    tag_values: List[str] = field(default_factory=list)
