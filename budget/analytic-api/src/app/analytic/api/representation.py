from typing import List

from pydantic import BaseModel


class BudgetExpenseAnalysisRequestRepresentation(BaseModel):
    year: int
    month: int
    tags: List[str]