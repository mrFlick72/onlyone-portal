import base64
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

    def find_by_id(revenue_id: RevenueId) -> Revenue:
        pass

    def find_by_data_range(
        user_name: UserName, start: Date, end: Date
    ) -> list[Revenue]:
        pass

    def save(revenue: Revenue) -> Revenue:
        pass

    def update(revenue: Revenue):
        pass

    def delete(id: RevenueId):
        pass
