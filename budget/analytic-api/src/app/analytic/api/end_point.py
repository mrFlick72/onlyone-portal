from typing import Annotated

from app.analytic.api.representation import (
    ReindexRequestRepresentation,
    ReindexResponseRepresentation,
    TagTotalRepresentation,
    TotalByTagAnalysisRequestRepresentation,
    TotalByYearAnalysisRequestRepresentation,
    YearTotalRepresentation,
)
from app.analytic.domain.reindex import ExpenseReindexService
from app.analytic.domain.service import BudgetExpenseAnalysisService
from app.container import ApplicationContainer
from fastapi import APIRouter, Depends
from dependency_injector.wiring import Provide, inject

analytic_end_point_router = APIRouter()


@analytic_end_point_router.put(
    "/api/analytic/budget/expense/total-by-tag"
)
@inject
async def budget_expense_total_by_tag(
    representation: TotalByTagAnalysisRequestRepresentation,
    analysis_service: Annotated[
        BudgetExpenseAnalysisService,
        Depends(
            Provide[
                ApplicationContainer.analytic_config_container.budget_expense_analysis_service
            ]
        ),
    ],
) -> list[TagTotalRepresentation]:
    totals = analysis_service.total_by_tag(
        representation.year, representation.month, representation.tags
    )

    return [
        TagTotalRepresentation(tag=total.tag, total=total.total.stringify_amount())
        for total in totals
    ]


@analytic_end_point_router.put(
    "/api/analytic/budget/expense/total-by-year"
)
@inject
async def budget_expense_total_by_year(
    representation: TotalByYearAnalysisRequestRepresentation,
    analysis_service: Annotated[
        BudgetExpenseAnalysisService,
        Depends(
            Provide[
                ApplicationContainer.analytic_config_container.budget_expense_analysis_service
            ]
        ),
    ],
) -> list[YearTotalRepresentation]:
    totals = analysis_service.total_by_year(
        representation.from_year, representation.to_year, representation.tag
    )

    return [
        YearTotalRepresentation(
            year=total.year.content, total=total.total.stringify_amount()
        )
        for total in totals
    ]


@analytic_end_point_router.post(
    "/api/analytic/budget/expense/reindex"
)
@inject
async def reindex_budget_expense(
    representation: ReindexRequestRepresentation,
    reindex_service: Annotated[
        ExpenseReindexService,
        Depends(
            Provide[ApplicationContainer.analytic_config_container.expense_reindex_service]
        ),
    ],
) -> ReindexResponseRepresentation:
    imported = reindex_service.reindex(
        representation.from_year, representation.to_year
    )
    return ReindexResponseRepresentation(imported=imported)
