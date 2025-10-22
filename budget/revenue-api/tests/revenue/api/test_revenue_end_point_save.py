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


def test_save_revenue_calls_service(client, mocker: MockerFixture):
    # Patch the converter used by the endpoint
    mocker.patch(
        "app.revenue.api.revenue_converter.from_representation_to_domain",
        return_value=Revenue(
            id=None,
            user_name=UserName("A_USER_NAME"),
            date=Date.iso_date_for("2018-10-10"),
            amount=Money.money_for("1.00"),
            note="A_NOTE",
        ),
        autospec=True,
    )

    # Mock the save service returned by the configuration provider
    mocked_save_revenue_use_case = mocker.Mock()

    mocker.patch(
        "app.revenue.config.RevenueConfigurationProvider.get_save_revenue_service",
        autospec=True,
        return_value=mocked_save_revenue_use_case,
    )

    response = client.post(
        "/budget/revenue",
        json={"date": "10/10/2018", "amount": "1.00", "note": "A_NOTE"},
    )

    assert response.status_code == 201

    mocked_save_revenue_use_case.save_revenue.assert_called_once_with(
        Revenue(
            id=None,
            user_name=UserName("A_USER_NAME"),
            date=Date.iso_date_for("2018-10-10"),
            amount=Money.money_for("1.00"),
            note="A_NOTE",
        )
    )
