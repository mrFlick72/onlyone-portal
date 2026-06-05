import json

from fastapi import APIRouter, Response

analytic_end_point_router = APIRouter()


@analytic_end_point_router.get("/api/analytic/hello", tags=["analytic"])
async def hello():
    return Response(
        status_code=200,
        content=json.dumps({"message": "hello world"}),
        media_type="application/json",
    )
