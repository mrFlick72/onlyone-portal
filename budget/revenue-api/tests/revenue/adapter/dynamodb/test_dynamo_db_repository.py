import boto3
import pytest
from botocore.exceptions import ClientError
from app.revenue.domain.revenue import Revenue, RevenueId
from app.time.domain.date import Date
from app.money.domain.money import Money
from app.user.domain.user import UserName

from app.revenue.adapter.dynamodb.dynamo_db_repository import DynamoDbRevenueIdProvider, DynamoDbRevenueRepository, UuidSaltGenerator


TABLE_NAME = "BUDGET_REVENUE"


@pytest.fixture(scope="module")
def dynamodb_resource():
    """Create a boto3 DynamoDB resource connected to LocalStack."""
    return boto3.resource(
        service_name="dynamodb",
        region_name="eu-central-1",
        aws_access_key_id="xxx",
        aws_secret_access_key="xxx",
        endpoint_url="http://localhost:4566",
    )


@pytest.fixture(scope="module")
def dynamodb_table(dynamodb_resource):
    """Create the BUDGET_REVENUE table before tests, delete after."""
    # Create table
    table = dynamodb_resource.create_table(
        TableName=TABLE_NAME,
        KeySchema=[
            {"AttributeName": "pk", "KeyType": "HASH"},  # partition key
            {"AttributeName": "range_key", "KeyType": "RANGE"},  # sort key
        ],
        AttributeDefinitions=[
            {"AttributeName": "pk", "AttributeType": "S"},
            {"AttributeName": "range_key", "AttributeType": "S"},
        ],
        BillingMode="PAY_PER_REQUEST",  # on-demand billing
    )

    # Wait for table to be created
    table.wait_until_exists()

    yield table

    # Teardown: delete table after tests
    try:
        table.delete()
        table.wait_until_not_exists()
    except ClientError as e:
        if "ResourceNotFoundException" not in str(e):
            raise


def test_find_by_data_range():
    pass


def test_save(dynamodb_table, dynamodb_resource):
    uut = DynamoDbRevenueRepository(
        dynamodb_resource, TABLE_NAME, FakeDynamoDbRevenueIdProvider()
    )

    revenue = Revenue(
        id=None,
        user_name=UserName("USER"),
        amount=Money.money_for("1.00"),
        date=Date.date_for("12/02/2000"),
        note="A_NOTE",
    )
    expected_revenue = Revenue(
        id=RevenueId("123-456"),
        user_name=UserName("USER"),
        amount=Money.money_for("1.00"),
        date=Date.date_for("12/02/2000"),
        note="A_NOTE",
    )

    uut.save(revenue)

    actual = uut.find_by_id(RevenueId("123-456"))

    assert expected_revenue == actual



def test_find_by_data_range(dynamodb_table, dynamodb_resource):
    uut = DynamoDbRevenueRepository(
        dynamodb_resource, TABLE_NAME, FakeDynamoDbRevenueIdProvider()
    )

    expected_revenue = [
        Revenue(
            id=RevenueId("123-4560"),
            user_name=UserName("USER"),
            amount=Money.money_for("17.50"),
            date=Date.date_for("06/01/2018"),
            note="Lanch",
        ),
        Revenue(
            id=RevenueId("123-4561"),
            user_name=UserName("USER"),
            amount=Money.money_for("10.50"),
            date=Date.date_for("12/02/2018"),
            note="Super Market",
        ),
        Revenue(
            id=RevenueId("123-4562"),
            user_name=UserName("USER"),
            amount=Money.money_for("17.50"),
            date=Date.date_for("13/02/2018"),
            note="Dinner",
        ),
        Revenue(
            id=RevenueId("123-4563"),
            user_name=UserName("USER"),
            amount=Money.money_for("17.50"),
            date=Date.date_for("22/02/2018"),
            note="Super Market",
        ),
    ]

    for revenue in expected_revenue:
        uut.save(revenue)

    actual = uut.find_by_data_range(
        UserName("USER"), Date.date_for("01/01/2018"), Date.date_for("31/12/2018")
    )
    
    assert expected_revenue == actual



def test_delete():
    pass


class FakeDynamoDbRevenueIdProvider(DynamoDbRevenueIdProvider):
    
    def __init__(self):
        super().__init__(UuidSaltGenerator())
    
    def generate_id(self, revenue: Revenue) -> RevenueId:
        if revenue.id is None:
            return RevenueId("123-456")
        return revenue.id
    

    def partition_key_from(self, date: Date, user_name: UserName) -> str:
        return "123"

    def range_key_from(self, date: Date) -> str:
        pass     
