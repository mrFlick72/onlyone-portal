from app.revenue.adapter.dynamodb.dynamo_db_repository import UuidSaltGenerator
from pytest_mock import MockerFixture

def test_uuid_revenue_id_salt_generator_happy_path(mocker: MockerFixture):
    uut = UuidSaltGenerator()
    
    mocker.patch("uuid.uuid4", return_value="123e4567-e89b-12d3-a456-426614174000")
    
    revenue_id = uut.new_salt()
    assert isinstance(revenue_id.content, str)
    assert revenue_id.content == "123e4567-e89b-12d3-a456-426614174000"  # UUID string length with hyphens