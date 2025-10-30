from pydantic import BaseModel
from app.time.domain.year import Year


class RevenueRepresentation(BaseModel):
    date: str
    amount: str
    note: str


class QueryParamRepresentation(BaseModel):
    year: Year
