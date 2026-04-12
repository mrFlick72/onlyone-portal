package revenue

import (
	"context"

	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/stretchr/testify/mock"
)

type RevenueRepositoryMock struct {
	mock.Mock
}

func (mock *RevenueRepositoryMock) FindFor(ctx context.Context, revenueId RevenueId) (*Revenue, error) {
	args := mock.Called(ctx, revenueId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Revenue), args.Error(1)
}

func (mock *RevenueRepositoryMock) FindByDateRange(ctx context.Context, start date.Date, end date.Date) ([]Revenue, error) {
	args := mock.Called(ctx, start, end)
	if args.Get(0) != nil {
		return args.Get(0).([]Revenue), args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *RevenueRepositoryMock) Save(ctx context.Context, revenue *Revenue) error {
	revenue.Id = "A_REVENUE_ID"
	args := mock.Called(ctx, revenue)
	return args.Error(0)
}

func (mock *RevenueRepositoryMock) Delete(ctx context.Context, revenueId RevenueId) error {
	args := mock.Called(ctx, revenueId)
	return args.Error(0)
}
