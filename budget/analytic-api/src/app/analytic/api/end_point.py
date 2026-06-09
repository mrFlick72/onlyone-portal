import json

from app.analytic.api.representation import BudgetExpenseAnalysisRequestRepresentation
from fastapi import APIRouter, Response
from dependency_injector.wiring import inject

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
):
    return Response(
        status_code=200,
        content=representation.model_dump_json(),
        media_type="application/json",
    )
