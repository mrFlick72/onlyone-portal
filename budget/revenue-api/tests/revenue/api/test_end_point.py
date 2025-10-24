import pytest
from app.server import app
from pytest_mock import MockerFixture
from app.user.domain.user import UserName
from app.time.domain.date import Date
from app.money.domain.money import Money
from app.revenue.domain.revenue import Revenue
from fastapi.testclient import TestClient

from app.revenue.domain.service import SaveRevenue
from app.revenue.api.converter import RevenueConverter
from app.revenue.api.representation import RevenueRepresentation
from app.revenue.domain.revenue import RevenueIdProvider


@pytest.fixture
def client():
    return TestClient(app)


def test_add_new_revenue(mocker: MockerFixture, client: TestClient):
    revenue = Revenue(
        id="generated-id",
        user_name=UserName("A_USER_NAME"),
        date=Date.iso_date_for("2018-10-10"),
        amount=Money.money_for("1.00"),
        note="A_NOTE",
    )
    revenue_representation = RevenueRepresentation(
        id="generated-id",
        date="10/10/2018",
        amount="1.00",
        note="A_NOTE",
    )

    mocked_converter = mocker.Mock(spec=RevenueConverter)
    mocked_converter.from_representation_to_domain.return_value = revenue

    mocked_save_revenue_use_case = mocker.Mock(spec=SaveRevenue)

    mocked_revenue_id_provider = mocker.Mock(spec=RevenueIdProvider)
    mocked_revenue_id_provider.generate_id.return_value = "generated-id"
    

    # Override dependencies
    app.container.revenue_config_container.revenue_converter.override(mocked_converter)
    app.container.revenue_config_container.save_revenue_service.override(
        mocked_save_revenue_use_case
    )
    app.container.revenue_config_container.revenue_id_provider.override(
        mocked_revenue_id_provider
    )

    response = client.post(
        "/budget/revenue",
        json={"date": "10/10/2018", "amount": "1.00", "note": "A_NOTE"},
    )

    assert response.status_code == 201
    mocked_converter.from_representation_to_domain.assert_called_once_with(
        revenue_representation
    )
    mocked_save_revenue_use_case.save.assert_called_once_with(revenue)
    mocked_revenue_id_provider.generate_id.assert_called_once()

    # cleanup
    app.container.revenue_config_container.revenue_converter.reset_override()
    app.container.revenue_config_container.save_revenue_service.reset_override()
    app.container.revenue_config_container.revenue_id_provider.reset_override()
