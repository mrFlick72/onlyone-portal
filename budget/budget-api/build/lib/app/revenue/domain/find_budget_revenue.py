from app.time.domain.date import Date
from app.time.domain.year import Year
from app.time.domain.month import Month

from app.user.domain.user_name_resolver import UserNameResolver
from app.revenue.domain.repository import RevenueRepository

class FindBudgetRevenue:

    def __init__(
        self, repository: RevenueRepository, userNameResolver: UserNameResolver
    ):
        self.repository = repository
        self.userNameResolver = userNameResolver

    def findBy(self, year: Year):
        return self.repository.find_by_data_range(
            self.userNameResolver.get_user_name(),
            Date.first_date_of_month(Month.JANUARY(), year),
            Date.last_date_of_month(Month.DECEMBER(), year),
        )
