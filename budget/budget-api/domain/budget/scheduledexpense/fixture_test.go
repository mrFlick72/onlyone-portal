package scheduledexpense

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ScheduledExpenseRepositoryMock struct {
	mock.Mock
}

func (m *ScheduledExpenseRepositoryMock) Save(ctx context.Context, scheduledExpense *ScheduledExpense) error {
	scheduledExpense.Id = "A_SCHEDULED_EXPENSE_ID"
	args := m.Called(ctx, scheduledExpense)
	return args.Error(0)
}

func (m *ScheduledExpenseRepositoryMock) FindAll(ctx context.Context) ([]ScheduledExpense, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ScheduledExpense), args.Error(1)
}
