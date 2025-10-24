from dependency_injector import containers, providers
from app.revenue.domain.service import SaveRevenue


class RevenueConfigContainer(containers.DeclarativeContainer):

    save_revenue_service = providers.Singleton(
        SaveRevenue,
    )
