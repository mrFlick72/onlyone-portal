from pydantic import BaseModel


class RevenueRepresentation(BaseModel):
    date: str
    amount: str
    note: str
