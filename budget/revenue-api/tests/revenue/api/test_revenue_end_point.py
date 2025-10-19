import pytest
from app.server import app
from pytest_mock import MockerFixture
from app.user.domain.user import UserName
from app.time.domain.date import Date
from app.money.domain.money import Money
from app.revenue.domain.revenue import Revenue


@pytest.fixture
def client():
    with app.test_client() as client:
        yield client


def test_add_new_revenue(client, mocker: MockerFixture):
    mocked_fromRepresentationToDomain = mocker.patch(
        "app.revenue.api.revenue_converter.fromRepresentationToDomain",
        return_value=Revenue(
            id=None,
            user_name=UserName("A_USER_NAME"),
            date=Date.iso_date_for("2018-10-10"),
            amount=Money.money_for("1.00"),
            note="A_NOTE",
        ),autospec=True
    )
    

    mocked_save_revenue_use_case = mocker.patch(
        "app.revenue.domain.service.SaveRevenue", autospec=True
    )

    mocked_save_revenue_use_case_factory = mocker.patch(
        "app.revenue.config.RevenueConfigurationProvider.get_save_revenue_service",
        autospec=True,
    )
    mocked_save_revenue_use_case_factory.return_value = mocked_save_revenue_use_case

    response = client.post(
        "/budget/revenue",
        json={"date": "10/10/2018", "amount": "1.00", "note": "A_NOTE"},
    )

    assert response.status_code == 201

    mocked_fromRepresentationToDomain.assert_called_once()
    mocked_save_revenue_use_case.save_revenue.assert_called_once_with(
        Revenue(
            id=None,
            user_name=UserName("A_USER_NAME"),
            date=Date.iso_date_for("2018-10-10"),
            amount=Money.money_for("1.00"),
            note="A_NOTE",
        )
    )
