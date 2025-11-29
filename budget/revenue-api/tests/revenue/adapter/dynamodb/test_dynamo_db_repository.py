import boto3
import pytest
from botocore.exceptions import ClientError
from app.revenue.domain.revenue import Revenue, RevenueId, RevenueIdProvider
from app.time.domain.date import Date
from app.money.domain.money import Money
from app.user.domain.user import UserName

from app.revenue.adapter.dynamodb.dynamo_db_repository import DynamoDbRevenueRepository


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
            {"AttributeName": "pk", "KeyType": "HASH"},   # partition key
            {"AttributeName": "sk", "KeyType": "RANGE"},  # sort key
        ],
        AttributeDefinitions=[
            {"AttributeName": "pk", "AttributeType": "S"},
            {"AttributeName": "sk", "AttributeType": "S"},
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


def test_find_by_id(dynamodb_table, dynamodb_resource):
    pass


def test_find_by_data_range():
    pass


def test_save(dynamodb_table, dynamodb_resource):
    uut = DynamoDbRevenueRepository(
        dynamodb_resource, TABLE_NAME, FixedDynamoDbRevenueIdProvider()
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


def test_update():
    pass


def test_delete():
    pass


class FixedDynamoDbRevenueIdProvider(RevenueIdProvider):
    def generate_id(self, revenue: Revenue) -> RevenueId:
        return RevenueId("123-456")
