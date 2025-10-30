from fastapi import APIRouter, Response, Depends
from dependency_injector.wiring import Provide, inject
from typing import Annotated

from app.revenue.api.representation import RevenueRequestRepresentation
from app.revenue.api.converter import RevenueConverter
from app.revenue.domain.service import SaveRevenue,DeleteRevenue
from app.container import ApplicationContainer
from app.revenue.domain.revenue import RevenueId

revenue_end_point_router = APIRouter()


@revenue_end_point_router.get("/budget/revenue")
async def get_revenue():
    return {}


@revenue_end_point_router.post("/budget/revenue")
@inject
async def save_revenue(
    representation: RevenueRequestRepresentation,
    save_revenue_service: Annotated[
        SaveRevenue,
        Depends(
            Provide[ApplicationContainer.revenue_config_container.save_revenue_service]
        ),
    ],
    converter: Annotated[
        RevenueConverter,
        Depends(
            Provide[ApplicationContainer.revenue_config_container.revenue_converter]
        ),
    ],
):
    revenue = converter.from_representation_to_domain(representation)
    save_revenue_service.save(revenue)
    return Response(status_code=201)


@revenue_end_point_router.put("/budget/revenue/{id}")
@inject
async def update_revenue(
    id: str,
    representation: RevenueRequestRepresentation,
    save_revenue_service: Annotated[
        SaveRevenue,
        Depends(
            Provide[ApplicationContainer.revenue_config_container.save_revenue_service]
        ),
    ],
    converter: Annotated[
        RevenueConverter,
        Depends(
            Provide[ApplicationContainer.revenue_config_container.revenue_converter]
        ),
    ],
):
    revenue = converter.from_representation_to_domain(representation)
    revenue.id = RevenueId(id)
    save_revenue_service.save(revenue)
    return Response(status_code=204)


@revenue_end_point_router.delete("/budget/revenue/{id}")
@inject
async def delete_revenue(    
    id: str,
    delete_revenue_service: Annotated[
        DeleteRevenue,
        Depends(
            Provide[ApplicationContainer.revenue_config_container.delete_revenue_service]
        ),
    ]):
    delete_revenue_service.delete(RevenueId(id))
    return Response(status_code=204)
