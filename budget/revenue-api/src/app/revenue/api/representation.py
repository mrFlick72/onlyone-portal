from typing import Optional
from pydantic import BaseModel


class RevenueRepresentation(BaseModel):
    id: Optional[str] = None
    date: str
    amount: str
    note: str
