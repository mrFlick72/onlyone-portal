from money.domain.money import Money
from time.domain.date import Date
from user.domain.user import UserName


class Revenue:

    def __init__(
        self, id, user_name: UserName, registration_date: Date, amount: Money, note: str
    ):
        self.id = id
        self.user_name = user_name
        self.registration_date = registration_date
        self.amount = amount
        self.note = note
