from dependency_injector import containers, providers
from app.user.container import UserConfigContainer
from app.revenue.container import RevenueConfigContainer


class ApplicationContainer(containers.DeclarativeContainer):

    user_config_container = providers.Container(
        UserConfigContainer,
    )

    revenue_config_container = providers.Container(
        RevenueConfigContainer,
        user_config_container=user_config_container,
    )

