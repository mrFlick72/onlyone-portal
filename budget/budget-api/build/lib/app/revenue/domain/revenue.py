from app.money.domain.money import Money
from app.time.domain.date import Date
from app.user.domain.user import UserName


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


