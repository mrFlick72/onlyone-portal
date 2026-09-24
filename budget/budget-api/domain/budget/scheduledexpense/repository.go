package scheduledexpense

import (
	"context"

	"github.com/mrflick72/budget/budget-api/domain/time/date"
)

// ScheduledExpenseRepository grows additively as tickets need more of it
// (#51 shipped Save+FindAll; #52 added FindFor; #53 added Delete; #54 adds
// UpdateStatus here), so this port never forces an unrelated ticket's
// implementation to exist first.
type ScheduledExpenseRepository interface {
	// Save persists scheduledExpense, generating its Id when empty.
	Save(ctx context.Context, scheduledExpense *ScheduledExpense) error

	// FindAll returns every ScheduledExpense owned by the user resolved from
	// ctx. There is no separate "owner" parameter — every other repository in
	// this service resolves the current user from ctx internally (see
	// DynamoDbRevenueRepository.FindByDateRange) rather than trusting a
	// caller-supplied value.
	FindAll(ctx context.Context) ([]ScheduledExpense, error)

	// FindFor returns the ScheduledExpense with the given id, scoped to the
	// user resolved from ctx — never the one on the struct passed to a
	// caller. Returns (nil, nil), not an error, when no such id exists in
	// that user's own partition: "not found" and "belongs to someone else"
	// are indistinguishable by construction here (PK = user_name), and
	// deliberately not treated as an error condition — callers (an Update
	// action's ownership check, a GET-by-id endpoint choosing 404 vs 500)
	// need to tell "not found" apart from a genuine failure.
	FindFor(ctx context.Context, id ScheduledExpenseId) (*ScheduledExpense, error)

	// Delete hard-deletes the definition row with the given id from the
	// current user's own partition. Returns ErrScheduledExpenseNotFound when
	// no such row exists there — including a row removed between a caller's
	// FindFor and this call.
	Delete(ctx context.Context, id ScheduledExpenseId) error

	// UpdateStatus sets only Status and LastEvaluatedDate on the row with the
	// given id in the current user's own partition, leaving every other field
	// as stored. Returns ErrScheduledExpenseNotFound when no such row exists
	// there.
	UpdateStatus(ctx context.Context, id ScheduledExpenseId, status Status, lastEvaluatedDate date.Date) error
}
