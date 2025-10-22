from dataclasses import dataclass
from app.money.domain.money import Money
from app.time.domain.date import Date
from app.user.domain.user import UserName


@dataclass
class RevenueId:
    content: str


class RevenueIdProvider:
    def generate_id(self) -> RevenueId:
        pass


class UuidRevenueIdProvider(RevenueIdProvider):
    def generate_id(self) -> RevenueId:
        import uuid

        return RevenueId(content=str(uuid.uuid4()))


@dataclass
class Revenue:
    id: RevenueId
    user_name: UserName
    date: Date
    amount: Money
    note: str
