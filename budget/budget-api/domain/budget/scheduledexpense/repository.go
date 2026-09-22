package scheduledexpense

import "context"

// ScheduledExpenseRepository is intentionally minimal for the create+list
// slice (#51) — Save and FindAll only. Update/Delete/status-transition
// methods are added additively by the tickets that need them (#52-#54), so
// this port grows without forcing an unrelated ticket's implementation to
// exist first.
type ScheduledExpenseRepository interface {
	// Save persists scheduledExpense, generating its Id when empty.
	Save(ctx context.Context, scheduledExpense *ScheduledExpense) error

	// FindAll returns every ScheduledExpense owned by the user resolved from
	// ctx. There is no separate "owner" parameter — every other repository in
	// this service resolves the current user from ctx internally (see
	// DynamoDbRevenueRepository.FindByDateRange) rather than trusting a
	// caller-supplied value.
	FindAll(ctx context.Context) ([]ScheduledExpense, error)
}
