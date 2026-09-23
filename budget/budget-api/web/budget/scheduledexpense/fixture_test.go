package scheduledexpense

import (
	"context"

	"github.com/gin-gonic/gin"
	domainscheduledexpense "github.com/mrflick72/budget/budget-api/domain/budget/scheduledexpense"
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

type ScheduledExpenseActionsMock struct {
	mock.Mock
}

func (m *ScheduledExpenseActionsMock) CreateScheduledExpense(ctx context.Context, se *domainscheduledexpense.ScheduledExpense) error {
	args := m.Called(ctx, se)
	return args.Error(0)
}

func (m *ScheduledExpenseActionsMock) FindScheduledExpenses(ctx context.Context) ([]domainscheduledexpense.ScheduledExpense, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domainscheduledexpense.ScheduledExpense), args.Error(1)
}

func (m *ScheduledExpenseActionsMock) FindScheduledExpense(ctx context.Context, id domainscheduledexpense.ScheduledExpenseId) (*domainscheduledexpense.ScheduledExpense, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainscheduledexpense.ScheduledExpense), args.Error(1)
}

func (m *ScheduledExpenseActionsMock) UpdateScheduledExpense(ctx context.Context, se *domainscheduledexpense.ScheduledExpense) error {
	args := m.Called(ctx, se)
	return args.Error(0)
}
