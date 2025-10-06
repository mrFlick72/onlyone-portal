from money.domain.money import Money
from time.domain.date import Date
from time.domain.month import Month
from user.domain.user import UserName
from user.domain.user_name_resolver import UserNameResolver
from time.domain.year import Year
from repository import RevenueRepository


class Revenue:

    def __init__(
        self, id, user_name: UserName, registration_date: Date, amount: Money, note: str
    ):
        self.id = id
        self.user_name = user_name
        self.registration_date = registration_date
        self.amount = amount
        self.note = note


class RevenueId:

    def __init__(self, content: str):
        self.content = content


class FindBudgetRevenue:

    def __init__(
        self, repository: RevenueRepository, userNameResolver: UserNameResolver
    ):
        self.repository = repository
        self.userNameResolver = userNameResolver

    def findBy(self, year: Year):
        self.repository.findByDateRange(
            self.userNameResolver.get_user_name().content,
            Date.firstDateOfMonth(Month.JANUARY, year),
            Date.lastDateOfMonth(Month.DECEMBER, year),
        )
