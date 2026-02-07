package web

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/stretchr/testify/mock"
)

type BudgetExpenseActionsMock struct {
	mock.Mock
}

func (mock *BudgetExpenseActionsMock) CreateBudgetExpense(ctx *context.Context, budgetExpense *expense.BudgetExpense) error {
	args := mock.Called(ctx, budgetExpense)
	return args.Error(0)
}

func (mock *BudgetExpenseActionsMock) UpdateBudgetExpense(ctx *context.Context, budgetExpense *expense.BudgetExpense) error {
	args := mock.Called(ctx, budgetExpense)
	return args.Error(0)
}

func (mock *BudgetExpenseActionsMock) FindSpentBudget(ctx *context.Context, month date.Month, year date.Year, searchTagKeys []tags.SearchTagKey) (*expense.SpentBudget, error) {
	args := mock.Called(ctx, month, year, searchTagKeys)
	return args.Get(0).(*expense.SpentBudget), args.Error(1)
}

func (mock *BudgetExpenseActionsMock) DeleteBudgetExpense(ctx *context.Context, id expense.BudgetExpenseId) error {
	args := mock.Called(ctx, id)
	return args.Error(0)
}	

type ContextFactoryConverterMock struct {
	mock.Mock
}

func (mock *ContextFactoryConverterMock) CreateContextFromGin(c *gin.Context) context.Context {
	args := mock.Called(c)
	return args.Get(0).(context.Context)
}	