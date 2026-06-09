from dataclasses import dataclass
from typing import List
from app.analytic.domain.expense import ExpenseRecord
from abc import ABC, abstractmethod


@dataclass
class BudgetExpenseAnalysisRequest:
    year: int
    month: int
    tags: List[str]


class ExpenseLoader(ABC):

    def __init__(self) -> None:
        super(self)

    @abstractmethod
    def expenseFor(request: BudgetExpenseAnalysisRequest) -> List[ExpenseRecord]: ...
