import os

from app.analytic.adapter.db.repository import PostgresExpenseProjectionRepository
from app.analytic.adapter.kafka.consumer import (
    BudgetExpenseConsumer,
    BudgetExpenseEventHandler,
)
from app.analytic.adapter.rest.source import RestBudgetExpenseSource
from app.analytic.domain.reindex import ExpenseReindexService
from app.analytic.domain.service import BudgetExpenseAnalysisService
from dependency_injector import containers, providers
from dotenv import load_dotenv


load_dotenv(dotenv_path=os.getenv("ANALYTIC_API_CONFIG_FILE_LOCATION"))


class AnalyticConfigContainer(containers.DeclarativeContainer):

    security_context_config_container = providers.DependenciesContainer()

    expense_projection_repository = providers.Singleton(
        PostgresExpenseProjectionRepository,
        database_url=os.getenv("DATABASE_URL"),
    )

    budget_expense_analysis_service = providers.Singleton(
        BudgetExpenseAnalysisService,
        repository=expense_projection_repository,
        security_context_resolver=security_context_config_container.security_context_resolver,
    )

    budget_expense_source = providers.Singleton(
        RestBudgetExpenseSource,
        budget_api_base_url=os.getenv("BUDGET_API_BASE_URL"),
        security_context_resolver=security_context_config_container.security_context_resolver,
    )

    expense_reindex_service = providers.Singleton(
        ExpenseReindexService,
        source=budget_expense_source,
        repository=expense_projection_repository,
    )

    budget_expense_event_handler = providers.Singleton(
        BudgetExpenseEventHandler,
        repository=expense_projection_repository,
    )

    budget_expense_consumer = providers.Singleton(
        BudgetExpenseConsumer,
        handler=budget_expense_event_handler,
        bootstrap_servers=os.getenv("KAFKA_BOOTSTRAP_SERVERS"),
        topic=os.getenv("KAFKA_EXPENSE_TOPIC", "budget-api.expense"),
        group_id=os.getenv("KAFKA_CONSUMER_GROUP", "analytic-api"),
    )
