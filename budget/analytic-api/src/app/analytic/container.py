import os

from app.analytic.adapter.service import RestExpenseLoader
from app.analytic.domain.service import BudgetExpenseAnalysisService
from dependency_injector import containers, providers
from dotenv import load_dotenv


load_dotenv(dotenv_path=os.getenv("ANALYTIC_API_CONFIG_FILE_LOCATION"))


class AnalyticConfigContainer(containers.DeclarativeContainer):

    security_context_config_container = providers.DependenciesContainer()

    expense_loader = providers.Singleton(
        RestExpenseLoader,
        security_context_resolver=security_context_config_container.security_context_resolver,
        budget_api_base_url=os.getenv("BUDGET_API_BASE_URL")
    )

    budget_expense_analysis_service = providers.Singleton(
        BudgetExpenseAnalysisService,
        expense_loader=expense_loader,
    )