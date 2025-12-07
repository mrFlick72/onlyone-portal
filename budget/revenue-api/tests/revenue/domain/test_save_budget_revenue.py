from pytest_mock import MockerFixture

from app.revenue.domain.service import SaveRevenue
from app.user.domain.user import UserName
from app.money.domain.money import Money
from app.revenue.domain.revenue import Revenue
from app.time.domain.date import Date

from app.revenue.domain.repository import RevenueRepository


def test_create_new_revenue_happy_path(mocker: MockerFixture):
    revenue = Revenue(
        id=None,
        user_name=UserName("A_USER_NAME"),
        date=Date.iso_date_for("2018-10-10"),
        amount=Money.money_for("1.00"),
        note="A_NOTE",
    )

    mocked_revenue_repository = mocker.Mock(spec=RevenueRepository)
    mocked_revenue_repository.save.return_value = None


    uut = SaveRevenue(
        repository=mocked_revenue_repository,
    )

    uut.save(revenue)

    mocked_revenue_repository.save.assert_called_once_with(revenue)


def test_update_revenue_happy_path(mocker: MockerFixture):
    revenue = Revenue(
        id="generated-id",
        user_name=UserName("A_USER_NAME"),
        date=Date.iso_date_for("2018-10-10"),
        amount=Money.money_for("1.00"),
        note="A_NOTE",
    )

    mocked_revenue_repository = mocker.Mock(spec=RevenueRepository)
    mocked_revenue_repository.save.return_value = None


    uut = SaveRevenue(
        repository=mocked_revenue_repository,
    )

    uut.save(revenue)

    mocked_revenue_repository.save.assert_called_once_with(revenue)