from dataclasses import dataclass, field
from datetime import date
from typing import List


@dataclass
class ExpenseRecord:
    id: str
    date: date
    amount: float
    note: str
    tag_values: List[str] = field(default_factory=list)
