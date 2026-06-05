from fastapi import APIRouter, Response, status

health_end_point_router = APIRouter()


@health_end_point_router.get("/health", tags=["management"])
async def health():
    return Response(status_code=status.HTTP_200_OK)
