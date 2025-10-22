from app.revenue.domain.revenue import UuidRevenueIdProvider
from pytest_mock import MockerFixture

def test_uuid_revenue_id_provider_happy_path(mocker: MockerFixture):
    uut = UuidRevenueIdProvider()
    
    mocker.patch("uuid.uuid4", return_value="123e4567-e89b-12d3-a456-426614174000")
    
    revenue_id = uut.generate_id()
    assert isinstance(revenue_id.content, str)
    assert revenue_id.content == "123e4567-e89b-12d3-a456-426614174000"  # UUID string length with hyphens