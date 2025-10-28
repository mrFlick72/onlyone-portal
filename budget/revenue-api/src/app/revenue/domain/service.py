from app.time.domain.date import Date
from app.time.domain.year import Year
from app.time.domain.month import Month

from app.user.domain.user_name_resolver import UserNameResolver
from app.revenue.domain.repository import RevenueRepository
from app.revenue.domain.revenue import Revenue, RevenueIdProvider


class SaveRevenue:

    def __init__(self, 
                    repository: RevenueRepository,
                    revenue_id_provider: RevenueIdProvider):           
            self.repository = repository
            self.revenue_id_provider = revenue_id_provider
        
    def save(self, revenue: Revenue):
        if revenue.id is None:
            revenue.id = self.revenue_id_provider.generate_id()
        self.repository.save(revenue)       



class FindBudgetRevenue:

    def __init__(
        self, repository: RevenueRepository, userNameResolver: UserNameResolver
    ):
        self.repository = repository
        self.userNameResolver = userNameResolver

    def findBy(self, year: Year):
        first_day_of_the_month = Date.first_date_of_month(Month.JANUARY(), year)
        last_day_of_the_month = Date.last_date_of_month(Month.DECEMBER(), year)
        current_user_name = self.userNameResolver.get_user_name()
        
        return self.repository.find_by_data_range(
            current_user_name,
            first_day_of_the_month,
            last_day_of_the_month,
        )
