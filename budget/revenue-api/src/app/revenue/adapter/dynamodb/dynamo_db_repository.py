from app.revenue.domain.repository import RevenueRepository
from app.revenue.domain.revenue import Revenue, RevenueId

from app.user.domain.user import UserName
from app.time.domain.date import Date


class DynamoDyIdGenerator:
    def budgetIdFrom(revenue: Revenue) -> RevenueId:
        pass

    def partition_key_from(date: Date, user_name: UserName) -> str:
        pass

    def partition_key_from(id: RevenueId) -> str:
        pass

    def range_key_from(id: RevenueId) -> str:
        pass


class DynamoDbRevenueRepository(RevenueRepository):

    def find_by_id(revenue_id: RevenueId) -> Revenue:
        pass

    def find_by_data_range(
        user_name: UserName, start: Date, end: Date
    ) -> list[Revenue]:
        pass

    def save(revenue: Revenue) -> Revenue:
        pass

    def update(revenue: Revenue):
        pass

    def delete(id: RevenueId):
        pass
