from dependency_injector import containers, providers
from app.revenue.domain.service import SaveRevenue
from app.revenue.api.converter import RevenueConverter
from app.revenue.domain.revenue import UuidRevenueIdProvider

class RevenueConfigContainer(containers.DeclarativeContainer):

    user_config_container = providers.DependenciesContainer()

    save_revenue_service = providers.Singleton(
        SaveRevenue,
    )

    revenue_converter = providers.Singleton(
        RevenueConverter,
        user_config_container.user_name_resolver,
    )
    
    revenue_id_provider = providers.Singleton(
        UuidRevenueIdProvider,
    )