from budget.revenue.domain.date import Date
from budget.revenue.domain.model.revenue import Revenue
from budget.user.domain.model.user_name import UserName
from budget.money.domain.money import Money


def a_revenue() -> Revenue:
    return Revenue(
        id="generated-id",
        user_name=UserName("A_USER_NAME"),
        date=Date.iso_date_for("2018-10-10"),
        amount=Money.money_for("1.00"),
        note="A_NOTE",
    )
