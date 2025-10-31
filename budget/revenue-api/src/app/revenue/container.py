from dependency_injector import containers, providers
from app.revenue.domain.service import SaveRevenue, DeleteRevenue
from app.revenue.api.converter import (
    RevenueConverter,
    QueryParamRepresentationConverter,
)
from app.revenue.domain.revenue import UuidRevenueIdProvider
from app.revenue.domain.repository import RevenueRepository
from app.revenue.domain.service import FindRevenue


class RevenueConfigContainer(containers.DeclarativeContainer):

    user_config_container = providers.DependenciesContainer()

    revenue_repository = providers.Singleton(
        RevenueRepository,
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
        revenue_id_provider=revenue_id_provider,
    )

    delete_revenue_service = providers.Singleton(
        DeleteRevenue,
        repository=revenue_repository,
        user_name_resolver=user_config_container.user_name_resolver,
    )
