import pytest
from pytest_mock import MockerFixture

from app.revenue.domain.service import DeleteRevenue
from app.user.domain.user import UserName
from app.money.domain.money import Money
from app.revenue.domain.revenue import Revenue
from app.time.domain.date import Date

from app.revenue.domain.repository import RevenueRepository
from app.user.domain.user_name_resolver import UserNameResolver
from app.revenue.domain.revenue import RevenueId

@pytest.fixture
def revenue_fixture() -> Revenue:
    return Revenue(
        id="generated_id",
        user_name=UserName("A_USER_NAME"),
        date=Date.iso_date_for("2018-10-10"),
        amount=Money.money_for("1.00"),
        note="A_NOTE",
    )


def test_delete_a_revenue_happy_path(mocker: MockerFixture, revenue_fixture: Revenue):
    mocked_revenue_repository = mocker.Mock(spec=RevenueRepository)
    mocked_revenue_repository.find_by_id.return_value = revenue_fixture

    mocked_user_name_resolver = mocker.Mock(spec=UserNameResolver)
    mocked_user_name_resolver.get_user_name.return_value = UserName("A_USER_NAME")

    uut = DeleteRevenue(
        repository=mocked_revenue_repository,
        user_name_resolver=mocked_user_name_resolver,
    )

    uut.delete(RevenueId("generated_id"))

    mocked_revenue_repository.delete.assert_called_once_with(RevenueId("generated_id"))
    mocked_user_name_resolver.get_user_name.assert_called_once()


def test_delete_a_revenue_fails(mocker: MockerFixture, revenue_fixture: Revenue):

    mocked_revenue_repository = mocker.Mock(spec=RevenueRepository)
    mocked_revenue_repository.find_by_id.return_value = revenue_fixture

    mocked_user_name_resolver = mocker.Mock(spec=UserNameResolver)
    mocked_user_name_resolver.get_user_name.return_value = UserName("ANOTHER_USER_NAME")

    uut = DeleteRevenue(
        repository=mocked_revenue_repository,
        user_name_resolver=mocked_user_name_resolver,
    )

    with pytest.raises(PermissionError):
        uut.delete(RevenueId("generated_id"))

    mocked_revenue_repository.delete.assert_not_called()
    mocked_user_name_resolver.get_user_name.assert_called_once()
