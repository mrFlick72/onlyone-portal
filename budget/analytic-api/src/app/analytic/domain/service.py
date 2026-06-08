from typing import List
from app.analytic.domain.expense import ExpenseRecord
from abc import ABC, abstractmethod


class ExpenseLoader(ABC):

    def __init__(self):
        super(self)

    @abstractmethod
    def expenseFor(
        self, year: int, month: int, tags: List[str]
    ) -> List[ExpenseRecord]: ...
