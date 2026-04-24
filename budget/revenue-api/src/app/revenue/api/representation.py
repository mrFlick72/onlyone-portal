from pydantic import BaseModel
from app.time.domain.year import Year


class RevenueRequestRepresentation(BaseModel):
    date: str
    amount: str
    note: str


class RevenueResponseRepresentation(BaseModel):
    id: str
    date: str
    amount: str
    note: str



class QueryParamRepresentation(BaseModel):
    year: Year
