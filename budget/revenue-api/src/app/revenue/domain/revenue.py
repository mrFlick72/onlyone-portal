from dataclasses import dataclass
from app.money.domain.money import Money
from app.time.domain.date import Date
from app.user.domain.user import UserName


@dataclass
class RevenueId:
    content: str


@dataclass
class Revenue:
    id: RevenueId
    user_name: UserName
    date: Date
    amount: Money
    note: str
