package scheduledexpense

import "context"

// ScheduledExpenseActions grows additively as tickets need more of it (#51
// shipped Create+FindScheduledExpenses; #52 adds FindScheduledExpense+Update
// here; Delete/Pause/Resume land in #53/#54).
type ScheduledExpenseActions interface {
	CreateScheduledExpense(ctx context.Context, scheduledExpense *ScheduledExpense) error
	FindScheduledExpenses(ctx context.Context) ([]ScheduledExpense, error)
	FindScheduledExpense(ctx context.Context, id ScheduledExpenseId) (*ScheduledExpense, error)
	UpdateScheduledExpense(ctx context.Context, scheduledExpense *ScheduledExpense) error
}

type ScheduledExpenseActionsFacade struct {
	CreateScheduledExpenseAction *CreateScheduledExpense
	FindScheduledExpensesAction  *FindScheduledExpenses
	FindScheduledExpenseAction   *FindScheduledExpense
	UpdateScheduledExpenseAction *UpdateScheduledExpense
}

func (facade *ScheduledExpenseActionsFacade) CreateScheduledExpense(ctx context.Context, scheduledExpense *ScheduledExpense) error {
	return facade.CreateScheduledExpenseAction.Execute(ctx, scheduledExpense)
}

func (facade *ScheduledExpenseActionsFacade) FindScheduledExpenses(ctx context.Context) ([]ScheduledExpense, error) {
	return facade.FindScheduledExpensesAction.Execute(ctx)
}

func (facade *ScheduledExpenseActionsFacade) FindScheduledExpense(ctx context.Context, id ScheduledExpenseId) (*ScheduledExpense, error) {
	return facade.FindScheduledExpenseAction.Execute(ctx, id)
}

func (facade *ScheduledExpenseActionsFacade) UpdateScheduledExpense(ctx context.Context, scheduledExpense *ScheduledExpense) error {
	return facade.UpdateScheduledExpenseAction.Execute(ctx, scheduledExpense)
}
