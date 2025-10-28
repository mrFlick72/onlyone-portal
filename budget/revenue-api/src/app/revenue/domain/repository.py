from app.user.domain.user import UserName
from app.time.domain.date import Date
from app.revenue.domain.revenue import Revenue, RevenueId


class RevenueRepository:

    def find_by_id(revenue_id: RevenueId) -> Revenue:
        pass

    def find_by_data_range(
        user_name: UserName, start: Date, end: Date
    ) -> list[Revenue]:
        pass

    def save(budgetRevenue: Revenue) -> Revenue:
        pass

    def update(budgetRevenue: Revenue):
        pass

    def delete(id: RevenueId):
        pass
