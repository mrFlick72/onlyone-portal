import os

from app.analytic.adapter.service import RestExpenseLoader
from dependency_injector import containers, providers


class AnalyticConfigContainer(containers.DeclarativeContainer):

    security_context_config_container = providers.DependenciesContainer()

    expense_loader = providers.Singleton(
        RestExpenseLoader,
        security_context_resolver=security_context_config_container,
        budget_api_base_url=os.getenv("BUDGET_API_BASE_URL")
    )