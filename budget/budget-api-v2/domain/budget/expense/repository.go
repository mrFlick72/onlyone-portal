package expense

import (
	"context"

	"github.com/mrflick72/budget/budget-api/domain/time/date"
)

type BudgetExpenseRepository interface {
	FindFor(ctx *context.Context, budgetExpenseId BudgetExpenseId) (*BudgetExpense, error)

	FindByDateRange(ctx *context.Context, star date.Date, end date.Date, searchTags []string) (*[]BudgetExpense, error)

	Save(ctx *context.Context, budgetExpense *BudgetExpense) error

	Delete(ctx *context.Context, idBudgetExpense BudgetExpenseId) error
}
