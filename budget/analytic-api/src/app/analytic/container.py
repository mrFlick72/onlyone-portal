from dependency_injector import containers, providers
from app.user.container import UserConfigContainer


class AnalyticConfigContainer(containers.DeclarativeContainer):

    user_config_container = providers.DependenciesContainer()
