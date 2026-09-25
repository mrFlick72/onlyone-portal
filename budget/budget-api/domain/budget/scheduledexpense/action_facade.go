package scheduledexpense

import "context"

// ScheduledExpenseActions grows additively as tickets need more of it (#51
// shipped Create+FindScheduledExpenses; #52 added FindScheduledExpense+Update;
// #53 added Delete; #54 adds Pause/Resume here).
type ScheduledExpenseActions interface {
	CreateScheduledExpense(ctx context.Context, scheduledExpense *ScheduledExpense) error
	FindScheduledExpenses(ctx context.Context) ([]ScheduledExpense, error)
	FindScheduledExpense(ctx context.Context, id ScheduledExpenseId) (*ScheduledExpense, error)
	UpdateScheduledExpense(ctx context.Context, scheduledExpense *ScheduledExpense) error
	DeleteScheduledExpense(ctx context.Context, id ScheduledExpenseId) error
	PauseScheduledExpense(ctx context.Context, id ScheduledExpenseId) error
	ResumeScheduledExpense(ctx context.Context, id ScheduledExpenseId) error
}

type ScheduledExpenseActionsFacade struct {
	CreateScheduledExpenseAction *CreateScheduledExpense
	FindScheduledExpensesAction  *FindScheduledExpenses
	FindScheduledExpenseAction   *FindScheduledExpense
	UpdateScheduledExpenseAction *UpdateScheduledExpense
	DeleteScheduledExpenseAction *DeleteScheduledExpense
	PauseScheduledExpenseAction  *PauseScheduledExpense
	ResumeScheduledExpenseAction *ResumeScheduledExpense
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

func (facade *ScheduledExpenseActionsFacade) DeleteScheduledExpense(ctx context.Context, id ScheduledExpenseId) error {
	return facade.DeleteScheduledExpenseAction.Execute(ctx, id)
}

func (facade *ScheduledExpenseActionsFacade) PauseScheduledExpense(ctx context.Context, id ScheduledExpenseId) error {
	return facade.PauseScheduledExpenseAction.Execute(ctx, id)
}

func (facade *ScheduledExpenseActionsFacade) ResumeScheduledExpense(ctx context.Context, id ScheduledExpenseId) error {
	return facade.ResumeScheduledExpenseAction.Execute(ctx, id)
}
