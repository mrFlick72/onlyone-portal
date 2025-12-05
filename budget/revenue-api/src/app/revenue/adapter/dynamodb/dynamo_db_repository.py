import logging
import os
from typing import Set
import boto3
import base64
from botocore.exceptions import ClientError

from app.money.domain.money import Money
from app.revenue.domain.repository import RevenueRepository
from app.revenue.domain.revenue import Revenue, RevenueId, RevenueIdProvider

from app.user.domain.user import UserName
from app.time.domain.date import Date
from itertools import chain


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
                f"{self.partition_key_from(date=date, user_name=user_name)}-{self.range_key_from(date=date)}"
            )
        else:
            return revenue.id

    def partition_key_from(self, date: Date, user_name: UserName) -> str:
        year = date.content.year
        user_name_content = user_name.content

        partition_key = f"{year}_{user_name_content}"
        utf_8_partition_key = partition_key.encode("utf-8")

        return base64.b64encode(utf_8_partition_key).decode("ascii")

    def range_key_from(self, date: Date) -> str:
        day_of_the_month = date.content.day
        month = date.content.month
        salt = self.salt_generator.new_salt()

        range_key = f"{month}_{day_of_the_month}_{salt}"
        utf_8_range_key = range_key.encode("utf-8")

        return base64.b64encode(utf_8_range_key).decode("ascii")


class DynamoDbRevenueRepository(RevenueRepository):

    ___logger = logging.getLogger(__name__)

    def __init__(
        self, dynamodb, table_name: str, id_generator: DynamoDbRevenueIdProvider
    ):
        super().__init__()
        self.dynamodb = dynamodb
        self.table_name = table_name
        self.id_generator = id_generator

    def find_by_id(self, revenue_id: RevenueId) -> Revenue:
        table = self.dynamodb.Table(self.table_name)
        try:
            response = table.get_item(
                Key={
                    "pk": revenue_id.content.split("-")[0],
                    "range_key": revenue_id.content.split("-")[1],
                }
            )
        except ClientError as e:
            self.___logger.error(e.response["Error"]["Message"])
        else:
            item = response.get("Item")
            if item:
                return Revenue(
                    id=RevenueId(item["budget_id"]),
                    user_name=UserName(item["user_name"]),
                    amount=Money.money_for(item["amount"]),
                    date=Date.iso_date_for(item["transaction_date"]),
                    note=item["note"],
                )
            else:
                return None

    def find_by_data_range(
        self, user_name: UserName, start: Date, end: Date
    ) -> list[Revenue]:        
        table = self.dynamodb.Table(self.table_name)
        partitions_keys = self.__primary_keys_for(user_name, start, end)
        revenues = list()
        for partitions_key in partitions_keys:
            response = table.query(
                KeyConditionExpression="pk = :pk",
                ExpressionAttributeValues={
                    ":pk": partitions_key,
                    ":start_date": start.iso_formatted_date(),
                    ":end_date": end.iso_formatted_date(),
                },
                FilterExpression="transaction_date >= :start_date AND transaction_date <= :end_date",
            )
            items = response.get("Items")
            revenues.append(
                [
                    Revenue(
                        id=RevenueId(item["budget_id"]),
                        user_name=UserName(item["user_name"]),
                        amount=Money.money_for(item["amount"]),
                        date=Date.iso_date_for(item["transaction_date"]),
                        note=item["note"],
                    )
                    for item in items
                ]
            )

            return list(chain.from_iterable(revenues))

    def __primary_keys_for(
        self, user_name: UserName, start: Date, end: Date
    ) -> Set[str]:
        primary_keys = set()
        current = start

        while current.content <= end.content:
            primary_keys.add(
                self.id_generator.partition_key_from(date=current, user_name=user_name)
            )
            current = current.plusYears(1)

        return primary_keys

    def save(self, revenue: Revenue) -> Revenue:
        revenue.id = self.id_generator.generate_id(revenue)
        table = self.dynamodb.Table(self.table_name)

        try:
            response = table.put_item(
                Item=self.__revenueAsDynamoDbItem(revenue),
            )
            self.___logger.debug("Item inserted successfully:", response)

        except ClientError as e:
            # If conditional insert fails or another error occurs
            self.___logger.error("Error inserting item:", e.response["Error"]["Message"])

    def update(self, revenue: Revenue):
        self.save(revenue)

    def delete(self, id: RevenueId):
        table = self.dynamodb.Table(self.table_name)
        
        table.delete_item(
            Key={
                "pk": id.content.split("-")[0],
                "range_key": id.content.split("-")[1],
            }
        )


    def __revenueAsDynamoDbItem(self, revenue: Revenue) -> dict:
        item = {
            "pk": revenue.id.content.split("-")[0],
            "range_key": revenue.id.content.split("-")[1],
            "budget_id": revenue.id.content,
            "user_name": revenue.user_name.content,
            "amount": revenue.amount.stringify_amount(),
            "transaction_date": revenue.date.content.isoformat(),
            "note": revenue.note,
        }
        return item
