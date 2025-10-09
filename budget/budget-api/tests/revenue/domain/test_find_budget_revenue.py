from pytest_mock import MockerFixture
from app.time.domain.date import Date
from app.time.domain.year import Year
from app.time.domain.month import Month
from app.money.domain.money import Money
from app.revenue.domain.revenue import Revenue, RevenueId
from app.revenue.domain.find_budget_revenue import FindBudgetRevenue


from app.user.domain.user import UserName


def test_find_yearly_revenue(mocker: MockerFixture):

    mocked_repository_response = [
        Revenue(
            id=RevenueId("A_REVENUE_ID"),
            user_name=UserName("A_USER_NAME"),
            registration_date=Date.iso_date_for("2021-02-02"),
            amount=Money.money_for("20.10"),
            note="A_NOTE",
        )
    ]
    mocked_repository = mocker.Mock()
    mocked_repository.find_by_data_range.return_value = mocked_repository_response

    mocked_user_name_resolver_response = UserName("A_USER_NAME")
    mocked_user_name_resolver = mocker.Mock()
    mocked_user_name_resolver.get_user_name.return_value = (
        mocked_user_name_resolver_response
    )

    uut = FindBudgetRevenue(mocked_repository, mocked_user_name_resolver)

    actual = uut.findBy(Year(2021))

    assert actual == mocked_repository_response
    mocked_repository.find_by_data_range.assert_called_once_with(
        UserName("A_USER_NAME"),
        Date.first_date_of_month(Month.JANUARY(), Year(2021)),
        Date.last_date_of_month(Month.DECEMBER(), Year(2021)),
    )
