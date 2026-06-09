import json
from typing import Annotated

from app.analytic.adapter.service import RestExpenseLoader
from app.analytic.api.representation import BudgetExpenseAnalysisRequestRepresentation
from app.container import ApplicationContainer
from fastapi import APIRouter, Depends, Response
from dependency_injector.wiring import Provide, inject

analytic_end_point_router = APIRouter()


@analytic_end_point_router.get("/api/analytic/hello", tags=["analytic"])
async def hello() -> Response:
    return Response(
        status_code=200,
        content=json.dumps({"message": "hello world"}),
        media_type="application/json",
    )


@analytic_end_point_router.put("/api/analytic/budget/expense", tags=["analytic"])
@inject
async def budget_analysis_for(
    representation: BudgetExpenseAnalysisRequestRepresentation,
    expense_loader: Annotated[
        RestExpenseLoader,
        Depends(
            Provide[ApplicationContainer.analytic_config_container.expense_loader]
        ),
    ],):
    print(expense_loader)
    return Response(
        status_code=200,
        content=representation.model_dump_json(),
        media_type="application/json",
    )
