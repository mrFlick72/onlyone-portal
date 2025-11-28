import boto3
from app.revenue.domain.revenue import Revenue, RevenueId
from app.time.domain.date import Date
from app.time.domain.year import Year
from app.time.domain.month import Month
from app.money.domain.money import Money
from app.user.domain.user import UserName

from app.revenue.adapter.dynamodb.dynamo_db_repository import DynamoDbRevenueRepository

def test_find_by_id():
    pass


def test_find_by_data_range():
    pass


def test_save():
    table_name = "BUDGET_REVENUE"
    dynamo = boto3.resource(
        service_name="dynamodb",
        region_name="eu-central-1",
        aws_access_key_id="xxx",
        aws_secret_access_key="xxx",
        endpoint_url="http://localhost:4566",
    )

    uut = DynamoDbRevenueRepository(dynamo, table_name)

    revenue = Revenue(
        id=None,
        user_name=UserName("USER"),
        amount=Money.money_for("1.00"),
        date=Date.date_for("12/02/2000"),
        note="A_NOTE",
    )
    expected_revenue = Revenue(
        id=RevenueId("123-465"),
        user_name=UserName("USER"),
        amount=Money.money_for("1.00"),
        date=Date.date_for("12/02/2000"),
        note="A_NOTE",
    )

    uut.save(revenue)

    actual = uut.find_by_id("123-456")

    assert expected_revenue == actual


def test_update():
    pass


def test_delete():
    pass
