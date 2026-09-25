package scheduledexpense

import (
	"context"

	"github.com/mrflick72/budget/budget-api/domain/time/date"
)

// ScheduledExpenseRepository grows additively as tickets need more of it
// (#51 shipped Save+FindAll; #52 added FindFor; #53 added Delete; #54 added
// UpdateStatus; #55 adds the generation methods), so this port never forces an
// unrelated ticket's implementation to exist first.
//
// It serves two callers. The user-facing methods (Save through UpdateStatus)
// are all scoped to the user resolved from ctx. The generation methods at the
// bottom exist for the scheduled expense job (ScheduledExpenseJob), which
// runs with no request in flight: FindAllActive is the one method that reads
// across every user and must never be called from a user-facing action.
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

	// --- ScheduledExpenseJob only (#55, ADR 0005) ---

	// FindAllActive returns EVERY user's ACTIVE definitions — the only method
	// not scoped to the user resolved from ctx, so reserved for the generation
	// engine; never call it from a user-facing action. Tags carry the keys and
	// the names stored at the definition's last save — never resolved through
	// tag-api, which needs a user access token the engine doesn't have.
	FindAllActive(ctx context.Context) ([]ScheduledExpense, error)

	// AdvanceLastEvaluatedDate sets LastEvaluatedDate on the definition with
	// the given id owned by the user resolved from ctx (the engine passes an
	// owner-only ctx). Returns ErrScheduledExpenseNotFound when that
	// definition no longer exists or is no longer ACTIVE (deleted or paused
	// mid-run).
	AdvanceLastEvaluatedDate(ctx context.Context, id ScheduledExpenseId, lastEvaluatedDate date.Date) error
}
