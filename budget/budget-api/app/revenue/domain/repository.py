from user.domain.user import UserName
from time.domain.date import Date
from revenue.domain.revenue import Revenue, RevenueId


class RevenueRepository:

    def findByDateRange(user_name: UserName, start: Date, end: Date) -> list[Revenue]:
        pass

    def save(budgetRevenue: Revenue) -> Revenue:
        pass

    def update(budgetRevenue: Revenue):
        pass

    def delete(id: RevenueId):
        pass
