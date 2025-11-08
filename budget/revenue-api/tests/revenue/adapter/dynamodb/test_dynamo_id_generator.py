from app.revenue.adapter.dynamodb.dynamo_db_repository import DynamoDbRevenueIdProvider
from app.revenue.domain.revenue import Revenue, RevenueId
from app.user.domain.user import UserName
from app.time.domain.date import Date
from app.revenue.adapter.dynamodb.dynamo_db_repository import SaltGenerator


def test_generate_id_for_empty_revenue_id():
    uut = DynamoDbRevenueIdProvider(FixedSaltGenerator())
    actual = uut.generate_id(
        Revenue(
            id=None,
            user_name=UserName("USER"),
            date=Date.date_for("06/01/2018"),
            amount=None,
            note="irrelevant",
        )
    )
    assert actual == RevenueId("MjAxOF9VU0VS-MV82X0FfU0FMVA==")

def test_generate_id_for_NOT_empty_revenue_id():
    uut = DynamoDbRevenueIdProvider(FixedSaltGenerator())
    actual = uut.generate_id(
        Revenue(
            id=RevenueId("A_FIXED_ID"),
            user_name=UserName("USER"),
            date=Date.date_for("06/01/2018"),
            amount=None,
            note="irrelevant",
        )
    )
    assert actual == RevenueId("A_FIXED_ID")


class FixedSaltGenerator(SaltGenerator):

    def new_salt(self) -> str:
        return "A_SALT"
