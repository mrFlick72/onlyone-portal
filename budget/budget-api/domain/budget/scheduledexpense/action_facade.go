package scheduledexpense

import "context"

// ScheduledExpenseActions is minimal for #51 (create+list) — Update/Delete/
// Pause/Resume are added additively by the tickets that need them (#52-#54).
type ScheduledExpenseActions interface {
	CreateScheduledExpense(ctx context.Context, scheduledExpense *ScheduledExpense) error
	FindScheduledExpenses(ctx context.Context) ([]ScheduledExpense, error)
}

type ScheduledExpenseActionsFacade struct {
	CreateScheduledExpenseAction *CreateScheduledExpense
	FindScheduledExpensesAction  *FindScheduledExpenses
}

func (facade *ScheduledExpenseActionsFacade) CreateScheduledExpense(ctx context.Context, scheduledExpense *ScheduledExpense) error {
	return facade.CreateScheduledExpenseAction.Execute(ctx, scheduledExpense)
}

func (facade *ScheduledExpenseActionsFacade) FindScheduledExpenses(ctx context.Context) ([]ScheduledExpense, error) {
	return facade.FindScheduledExpensesAction.Execute(ctx)
}
