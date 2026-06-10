import json
from typing import Annotated

from app.analytic.adapter.service import RestExpenseLoader
from app.analytic.api.representation import (
    BudgetExpenseAnalysisRequestRepresentation,
    BudgetExpenseAnalysisResponseRepresentation,
    TagTotalRepresentation,
    TotalByTagAnalysisRequestRepresentation,
    TotalByYearAnalysisRequestRepresentation,
    YearTotalRepresentation,
)
from app.analytic.domain.expense import BudgetExpenseAnalysisRequest
from app.analytic.domain.service import BudgetExpenseAnalysisService
from app.container import ApplicationContainer
from fastapi import APIRouter, Depends, Response
from dependency_injector.wiring import Provide, inject

analytic_end_point_router = APIRouter()



@analytic_end_point_router.put("/api/analytic/budget/expense")
@inject
async def budget_analysis_for(
    representation: BudgetExpenseAnalysisRequestRepresentation,
    expense_loader: Annotated[
        RestExpenseLoader,
        Depends(Provide[ApplicationContainer.analytic_config_container.expense_loader]),
    ],
) -> list[BudgetExpenseAnalysisResponseRepresentation]:
    expenses = expense_loader.expenseFor(
        BudgetExpenseAnalysisRequest(
            representation.year, representation.month, representation.tags
        )
    )

    return [
        BudgetExpenseAnalysisResponseRepresentation(
            date=exp.date.formatted_date(),
            amount=exp.amount.stringify_amount(),
            tag_values=exp.tag_values,
        )
        for exp in expenses
    ]


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
