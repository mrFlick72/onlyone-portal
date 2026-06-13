package expense

import "context"

type BudgetExpenseEventPublisher interface {
	CreateBudgetExpense(ctx context.Context, expense BudgetExpense) error
	UpdateBudgetExpense(ctx context.Context, expense BudgetExpense) error
	DeleteBudgetExpense(ctx context.Context, expense BudgetExpense) error
}
