from dependency_injector import containers, providers


class AnalyticConfigContainer(containers.DeclarativeContainer):

    security_context_config_container = providers.DependenciesContainer()
