package expense

import (
	"context"

	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/stretchr/testify/mock"
)

type BudgetExpenseRepositoryMock struct {
	mock.Mock
}

func (mock *BudgetExpenseRepositoryMock) FindFor(ctx *context.Context, budgetExpenseId BudgetExpenseId) (*BudgetExpense, error) {
	args := mock.Called(ctx, budgetExpenseId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BudgetExpense), args.Error(1)
}

func (mock *BudgetExpenseRepositoryMock) FindByDateRange(ctx *context.Context,star date.Date, end date.Date, searchTags []string) (*[]BudgetExpense, error) {
	argrs := mock.Called(ctx, star, end, searchTags)
	if argrs.Get(0) != nil {
		return argrs.Get(0).(*[]BudgetExpense), argrs.Error(1)
	}
	return nil, nil
}
func (mock *BudgetExpenseRepositoryMock) Save(ctx *context.Context, budgetExpense *BudgetExpense) error {
	budgetExpense.Id = "A_BUDGET_ID"
	args := mock.Called(ctx, budgetExpense)
	return args.Error(0)

}
func (mock *BudgetExpenseRepositoryMock) Delete(ctx *context.Context, idBudgetExpense BudgetExpenseId) error {
	args := mock.Called(ctx, idBudgetExpense)
	return args.Error(0)
}

type SearchTagRepositoryMock struct {
	mock.Mock
}

func (mock *SearchTagRepositoryMock) GetTagBy(ctx *context.Context, searchTagKey tags.SearchTagKey) (*tags.SearchTag, error) {
	args := mock.Called(ctx, searchTagKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tags.SearchTag), args.Error(1)
}

