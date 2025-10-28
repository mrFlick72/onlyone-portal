from fastapi import APIRouter, Response, Depends
from dependency_injector.wiring import Provide, inject
from typing import Annotated

from app.revenue.api.representation import RevenueRepresentation
from app.revenue.api.converter import RevenueConverter
from app.revenue.domain.service import SaveRevenue
from app.container import ApplicationContainer

revenue_end_point_router = APIRouter()


@revenue_end_point_router.route("/budget/revenue", methods=["GET"])
async def get_revenue():
    return {}

@revenue_end_point_router.post("/budget/revenue")
@inject
async def save_revenue(
    representation: RevenueRepresentation,
    save_revenue_service: Annotated[
        SaveRevenue, Depends(Provide[ApplicationContainer.revenue_config_container.save_revenue_service])
    ],
    converter: Annotated[
        RevenueConverter, Depends(Provide[ApplicationContainer.revenue_config_container.revenue_converter])
    ]
):
    revenue = converter.from_representation_to_domain(representation)
    save_revenue_service.save(revenue)
    return Response(status_code=201)


# @revenue_end_point_router.route("/budget/revenue/{id}", methods=["PUT"])
# async def update_revenue():
#     return {}


# @revenue_end_point_router.route("/budget/revenue/{id}", methods=["DELETE"])
# async def delete_revenue():
#     return {}
