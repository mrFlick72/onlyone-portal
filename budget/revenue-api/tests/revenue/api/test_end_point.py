import json
import pytest
from app.server import app
from pytest_mock import MockerFixture
from app.user.domain.user import UserName
from app.time.domain.date import Date
from app.time.domain.year import Year
from app.money.domain.money import Money
from app.revenue.domain.revenue import Revenue, RevenueId
from fastapi.testclient import TestClient

from app.revenue.domain.service import SaveRevenue, DeleteRevenue, FindRevenue
from app.revenue.api.converter import (
    RevenueConverter,
    QueryParamRepresentationConverter,
)
from app.revenue.api.representation import (
    RevenueRequestRepresentation,
    RevenueResponseRepresentation,
    QueryParamRepresentation,
)


@pytest.fixture
def client():
    return TestClient(app)


def test_add_new_revenue(mocker: MockerFixture, client: TestClient):
    revenue = Revenue(
        id=None,
        user_name=UserName("A_USER_NAME"),
        date=Date.iso_date_for("2018-10-10"),
        amount=Money.money_for("1.00"),
        note="A_NOTE",
    )

    revenue_representation = RevenueRequestRepresentation(
        date="10/10/2018",
        amount="1.00",
        note="A_NOTE",
    )

    mocked_converter = mocker.Mock(spec=RevenueConverter)
    mocked_converter.from_representation_to_domain.return_value = revenue

    mocked_save_revenue_use_case = mocker.Mock(spec=SaveRevenue)

    # Override dependencies
    app.container.revenue_config_container.revenue_converter.override(mocked_converter)
    app.container.revenue_config_container.save_revenue_service.override(
        mocked_save_revenue_use_case
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

    # cleanup
    app.container.revenue_config_container.revenue_converter.reset_override()
    app.container.revenue_config_container.save_revenue_service.reset_override()


def test_read_revenue_by_year(mocker: MockerFixture, client: TestClient):
    mocked_find_revenue_use_case = mocker.Mock(spec=FindRevenue)
    mocked_find_revenue_use_case.find_by.return_value = [
        Revenue(
            id=RevenueId("revenue-id-1"),
            user_name=UserName("A_USER_NAME"),
            date=Date.iso_date_for("2023-03-15"),
            amount=Money.money_for("10.00"),
            note="A NOTE",
        ),
        Revenue(
            id=RevenueId("revenue-id-2"),
            user_name=UserName("A_USER_NAME"),
            date=Date.iso_date_for("2023-03-16"),
            amount=Money.money_for("15.00"),
            note="ANOTHER NOTE",
        ),
    ]

    mocked_converter = mocker.Mock(spec=QueryParamRepresentationConverter)
    mocked_converter.from_query_param_to_year.return_value = QueryParamRepresentation(
        year=Year(2023)
    )

    # Override dependencies
    app.container.revenue_config_container.query_param_converter.override(
        mocked_converter
    )
    app.container.revenue_config_container.find_revenue_service.override(
        mocked_find_revenue_use_case
    )

    response = client.get("/budget/revenue?q=year=2023")

    assert response.status_code == 200
    assert response.json() == [
        {
            "id": "revenue-id-1",
            "date": "15/03/2023",
            "amount": "10.00",
            "note": "A NOTE",
        },
        {
            "id": "revenue-id-2",
            "date": "16/03/2023",
            "amount": "15.00",
            "note": "ANOTHER NOTE",
        },
    ]

    mocked_converter.from_query_param_to_year.assert_called_once_with("year=2023")
    mocked_find_revenue_use_case.find_by.assert_called_once_with(Year(2023))


def test_update_revenue(mocker: MockerFixture, client: TestClient):
    converted_revenue = Revenue(
        id=None,
        user_name=UserName("A_USER_NAME"),
        date=Date.iso_date_for("2018-10-10"),
        amount=Money.money_for("1.00"),
        note="A_NOTE",
    )
    revenue = Revenue(
        id=RevenueId("123"),
        user_name=UserName("A_USER_NAME"),
        date=Date.iso_date_for("2018-10-10"),
        amount=Money.money_for("1.00"),
        note="A_NOTE",
    )

    revenue_representation = RevenueRequestRepresentation(
        date="10/10/2018",
        amount="1.00",
        note="A_NOTE",
    )

    mocked_converter = mocker.Mock(spec=RevenueConverter)
    mocked_converter.from_representation_to_domain.return_value = converted_revenue

    mocked_save_revenue_use_case = mocker.Mock(spec=SaveRevenue)

    # Override dependencies
    app.container.revenue_config_container.revenue_converter.override(mocked_converter)
    app.container.revenue_config_container.save_revenue_service.override(
        mocked_save_revenue_use_case
    )

    response = client.put(
        "/budget/revenue/123",
        json={"date": "10/10/2018", "amount": "1.00", "note": "A_NOTE"},
    )

    assert response.status_code == 204
    mocked_converter.from_representation_to_domain.assert_called_once_with(
        revenue_representation
    )
    mocked_save_revenue_use_case.save.assert_called_once_with(revenue)

    # cleanup
    app.container.revenue_config_container.revenue_converter.reset_override()
    app.container.revenue_config_container.save_revenue_service.reset_override()


def test_delete_revenue(mocker: MockerFixture, client: TestClient):
    mocked_delete_revenue_use_case = mocker.Mock(spec=DeleteRevenue)

    # Override dependencies
    app.container.revenue_config_container.delete_revenue_service.override(
        mocked_delete_revenue_use_case
    )

    response = client.delete("/budget/revenue/123")

    assert response.status_code == 204

    mocked_delete_revenue_use_case.delete.assert_called_once_with(RevenueId("123"))

    # cleanup
    app.container.revenue_config_container.delete_revenue_service.reset_override()
