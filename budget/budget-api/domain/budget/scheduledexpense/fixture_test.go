package scheduledexpense

import (
	"context"

	"github.com/mrflick72/budget/budget-api/domain/time/date"

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

func (m *ScheduledExpenseRepositoryMock) FindFor(ctx context.Context, id ScheduledExpenseId) (*ScheduledExpense, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ScheduledExpense), args.Error(1)
}

func (m *ScheduledExpenseRepositoryMock) Delete(ctx context.Context, id ScheduledExpenseId) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ScheduledExpenseRepositoryMock) UpdateStatus(ctx context.Context, id ScheduledExpenseId, status Status, lastEvaluatedDate date.Date) error {
	args := m.Called(ctx, id, status, lastEvaluatedDate)
	return args.Error(0)
}
