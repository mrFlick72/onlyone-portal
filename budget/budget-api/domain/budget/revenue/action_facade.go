package revenue

import (
	"context"

	"github.com/mrflick72/budget/budget-api/domain/time/date"
)

type RevenueActions interface {
	CreateRevenue(ctx context.Context, revenue *Revenue) error
	UpdateRevenue(ctx context.Context, revenue *Revenue) error
	FindRevenue(ctx context.Context, year date.Year) ([]Revenue, error)
	DeleteRevenue(ctx context.Context, id RevenueId) error
}

type RevenueActionsFacade struct {
	CreateRevenueAction *CreateRevenue
	UpdateRevenueAction *UpdateRevenue
	FindRevenueAction   *FindRevenue
	DeleteRevenueAction *DeleteRevenue
}

func (facade *RevenueActionsFacade) CreateRevenue(ctx context.Context, revenue *Revenue) error {
	return facade.CreateRevenueAction.Execute(ctx, revenue)
}

func (facade *RevenueActionsFacade) UpdateRevenue(ctx context.Context, revenue *Revenue) error {
	return facade.UpdateRevenueAction.Execute(ctx, revenue)
}

func (facade *RevenueActionsFacade) FindRevenue(ctx context.Context, year date.Year) ([]Revenue, error) {
	return facade.FindRevenueAction.Execute(ctx, year)
}

func (facade *RevenueActionsFacade) DeleteRevenue(ctx context.Context, id RevenueId) error {
	return facade.DeleteRevenueAction.Execute(ctx, id)
}
