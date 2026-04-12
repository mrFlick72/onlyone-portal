package revenue

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/budget/budget-api/domain/budget/revenue"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/stretchr/testify/mock"
)

func SetUpRouter() *gin.Engine {
	return gin.Default()
}

type ContextFactoryConverterMock struct {
	mock.Mock
}

func (m *ContextFactoryConverterMock) CreateContextFromGin(c *gin.Context) context.Context {
	args := m.Called(c)
	return args.Get(0).(context.Context)
}

type RevenueActionsMock struct {
	mock.Mock
}

func (m *RevenueActionsMock) CreateRevenue(ctx context.Context, r *revenue.Revenue) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *RevenueActionsMock) UpdateRevenue(ctx context.Context, r *revenue.Revenue) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *RevenueActionsMock) FindRevenue(ctx context.Context, year date.Year) ([]revenue.Revenue, error) {
	args := m.Called(ctx, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]revenue.Revenue), args.Error(1)
}

func (m *RevenueActionsMock) DeleteRevenue(ctx context.Context, id revenue.RevenueId) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
