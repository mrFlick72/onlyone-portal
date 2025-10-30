from dependency_injector import containers, providers
from app.revenue.domain.service import SaveRevenue, DeleteRevenue
from app.revenue.api.converter import RevenueConverter
from app.revenue.domain.revenue import UuidRevenueIdProvider
from app.revenue.domain.repository import RevenueRepository
from app.user.domain.user_name_resolver import UserNameResolver


class RevenueConfigContainer(containers.DeclarativeContainer):

    user_config_container = providers.DependenciesContainer()

    revenue_repository = providers.Singleton(
        RevenueRepository,
    )

    revenue_converter = providers.Singleton(
        RevenueConverter,
        user_config_container.user_name_resolver,
    )

    revenue_id_provider = providers.Singleton(
        UuidRevenueIdProvider,
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
