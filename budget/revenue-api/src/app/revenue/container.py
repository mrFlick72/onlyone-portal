import os
import boto3
from dependency_injector import containers, providers
from app.revenue.adapter.dynamodb.dynamo_db_repository import (
    DynamoDbRevenueIdProvider,
    DynamoDbRevenueRepository,
    UuidSaltGenerator,
)
from app.revenue.domain.service import SaveRevenue, DeleteRevenue
from app.revenue.api.converter import (
    RevenueConverter,
    QueryParamRepresentationConverter,
)
from app.revenue.domain.revenue import UuidRevenueIdProvider
from app.revenue.domain.service import FindRevenue

from dotenv import load_dotenv


load_dotenv(dotenv_path=os.getenv("BUDGET_API_CONFIG_FILE_LOCATION"))


class RevenueConfigContainer(containers.DeclarativeContainer):

    user_config_container = providers.DependenciesContainer()

    # Singleton DynamoDB resource
    dynamodb_resource = providers.Singleton(
        boto3.resource,
        service_name="dynamodb",
        aws_access_key_id=os.getenv("AWS_ACCESS_KEY_ID", "xxx"),
        aws_secret_access_key=os.getenv("AWS_SECRET_ACCESS_KEY", "xxx"),
        region_name=os.getenv("AWS_REGION", "eu-central-1"),
    )

    # Singleton ID provider
    salt_generator = providers.Singleton(UuidSaltGenerator)

    id_provider = providers.Singleton(
        DynamoDbRevenueIdProvider,
        salt_generator=salt_generator,
    )

    # Singleton repository that depends on the above
    revenue_repository = providers.Singleton(
        DynamoDbRevenueRepository,
        dynamodb=dynamodb_resource,
        table_name=os.getenv("REVENUE_DYNAMO_DB_TABLE_NAME", "BUDGET_REVENUE"),
        id_generator=id_provider,
    )

    revenue_converter = providers.Singleton(
        RevenueConverter,
        user_config_container.user_name_resolver,
    )

    query_param_converter = providers.Singleton(
        QueryParamRepresentationConverter,
    )

    revenue_id_provider = providers.Singleton(
        UuidRevenueIdProvider,
    )

    find_revenue_service = providers.Singleton(
        FindRevenue,
        repository=revenue_repository,
        user_name_resolver=user_config_container.user_name_resolver,
    )

    save_revenue_service = providers.Singleton(
        SaveRevenue,
        repository=revenue_repository,
    )

    delete_revenue_service = providers.Singleton(
        DeleteRevenue,
        repository=revenue_repository,
        user_name_resolver=user_config_container.user_name_resolver,
    )
