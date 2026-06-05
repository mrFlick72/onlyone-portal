from dependency_injector import containers, providers
from app.user.container import UserConfigContainer
from app.analytic.container import AnalyticConfigContainer


class ApplicationContainer(containers.DeclarativeContainer):

    user_config_container = providers.Container(
        UserConfigContainer,
    )

    analytic_config_container = providers.Container(
        AnalyticConfigContainer,
        user_config_container=user_config_container,
    )
