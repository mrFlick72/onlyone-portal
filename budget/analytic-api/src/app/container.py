from app.infrastructure.security.security_container import SecurityConfigContainer
from dependency_injector import containers, providers
from app.analytic.container import AnalyticConfigContainer


class ApplicationContainer(containers.DeclarativeContainer):

    security_context_config_container = providers.Container(
        SecurityConfigContainer,
    )

    analytic_config_container = providers.Container(
        AnalyticConfigContainer,
        security_context_config_container=security_context_config_container,
    )
