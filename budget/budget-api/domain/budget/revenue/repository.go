package revenue

import (
	"context"

	"github.com/mrflick72/budget/budget-api/domain/time/date"
)

type RevenueRepository interface {
	FindFor(ctx context.Context, revenueId RevenueId) (*Revenue, error)

	FindByDateRange(ctx context.Context, start date.Date, end date.Date) ([]Revenue, error)

	Save(ctx context.Context, revenue *Revenue) error

	Delete(ctx context.Context, revenueId RevenueId) error
}
