package domain

import "context"

type BudgetExpenseRepository interface {
	FindFor(ctx *context.Context, budgetExpenseId BudgetExpenseId) (*BudgetExpense, error)

	FindByDateRange(ctx *context.Context, userName UserName, star Date, end Date, searchTags []string) (*[]BudgetExpense, error)

	Save(ctx *context.Context, budgetExpense *BudgetExpense) error

	Delete(ctx *context.Context, idBudgetExpense BudgetExpenseId) error
}
