from pytest_mock import MockerFixture

from app.revenue.domain.service import SaveRevenue
from app.user.domain.user import UserName
from app.money.domain.money import Money
from app.revenue.domain.revenue import Revenue
from app.time.domain.date import Date

def test_happy_path(mocker: MockerFixture):
    revenue = Revenue(
        id="generated-id",
        user_name=UserName("A_USER_NAME"),
        date=Date.iso_date_for("2018-10-10"),
        amount=Money.money_for("1.00"),
        note="A_NOTE",
    )
    uut = SaveRevenue()

    uut.save(revenue)
