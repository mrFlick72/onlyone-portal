import boto3
import base64
from botocore.exceptions import ClientError

from app.revenue.domain.repository import RevenueRepository
from app.revenue.domain.revenue import Revenue, RevenueId, RevenueIdProvider

from app.user.domain.user import UserName
from app.time.domain.date import Date


class SaltGenerator:
    def new_salt(self) -> str:
        pass


class UuidSaltGenerator(SaltGenerator):

    def new_salt(self):
        import uuid

        return RevenueId(content=str(uuid.uuid4()))


class DynamoDbRevenueIdProvider(RevenueIdProvider):

    def __init__(self, salt_generator: SaltGenerator):
        super().__init__()
        self.salt_generator = salt_generator

    def generate_id(self, revenue: Revenue) -> RevenueId:
        if revenue.id is None:
            date = revenue.date
            user_name = revenue.user_name
            return RevenueId(
                f"{self.__partition_key_from(date=date, user_name=user_name)}-{self.__range_key_from(date=date)}"
            )
        else:
            return revenue.id

    def __partition_key_from(self, date: Date, user_name: UserName) -> str:
        year = date.content.year
        user_name_content = user_name.content

        partition_key = f"{year}_{user_name_content}"
        utf_8_partition_key = partition_key.encode("utf-8")

        return base64.b64encode(utf_8_partition_key).decode("ascii")

    def __range_key_from(self, date: Date) -> str:
        day_of_the_month = date.content.day
        month = date.content.month
        salt = self.salt_generator.new_salt()

        range_key = f"{month}_{day_of_the_month}_{salt}"
        utf_8_range_key = range_key.encode("utf-8")

        return base64.b64encode(utf_8_range_key).decode("ascii")


class DynamoDbRevenueRepository(RevenueRepository):

    def __init__(self, dynamodb, table_name: str):
        super().__init__()
        self.dynamodb = dynamodb
        self.table_name = table_name

    def find_by_id(self, revenue_id: RevenueId) -> Revenue:
        pass

    def find_by_data_range(
        self, user_name: UserName, start: Date, end: Date
    ) -> list[Revenue]:
        pass

    def save(self, revenue: Revenue) -> Revenue:
        print("Saving revenue to DynamoDB:", revenue)
        print("table_name: ", self.table_name)
        table = self.dynamodb.Table(self.table_name)
        try:
            # response = table.save(revenue)
            response = table.put_item(
                Item=self.__revenueAsDynamoDbItem(revenue),
                ConditionExpression="attribute_not_exists(pk)",  # optional: avoid overwriting
            )
            print("Item inserted successfully:", response)

        except ClientError as e:
            # If conditional insert fails or another error occurs
            print("Error inserting item:", e.response["Error"]["Message"])

    def update(self, revenue: Revenue):
        pass

    def delete(self, id: RevenueId):
        pass

    def __revenueAsDynamoDbItem(self, revenue: Revenue) -> dict:
        item = {
            "pk": revenue.id.content.split("-")[0],
            "sk": revenue.id.content.split("-")[1],
            "budget_id": revenue.id.content,
            "user_name": revenue.user_name.content,
            "amount": str(revenue.amount.stringify_amount()),
            "transaction_date": revenue.date.content.isoformat(),
            "note": revenue.note,
        }
        return item
