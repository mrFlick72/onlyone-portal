package scheduledexpense

import (
	"context"
	"errors"
	"sort"

	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

// applyDefaultTagIfMissing defaults a ScheduledExpense with no tags to the
// shared UNKNOWN sentinel — the same invariant CreateBudgetExpense and
// CreateRevenue enforce, so a template with no tags reads and generates
// exactly like an explicitly-untagged expense (see domain/tags.UnknownSentinel).
func applyDefaultTagIfMissing(scheduledExpense *ScheduledExpense) {
	if len(scheduledExpense.Tags) == 0 {
		scheduledExpense.Tags = []tags.SearchTag{tags.UnknownSentinel()}
	}
}

type CreateScheduledExpense struct {
	Repository ScheduledExpenseRepository
}

func (action *CreateScheduledExpense) Execute(ctx context.Context, scheduledExpense *ScheduledExpense) error {
	applyDefaultTagIfMissing(scheduledExpense)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}
	scheduledExpense.UserName = *user.UserName
	scheduledExpense.Status = StatusActive
	return action.Repository.Save(ctx, scheduledExpense)
}

type FindScheduledExpense struct {
	Repository ScheduledExpenseRepository
}

func (action *FindScheduledExpense) Execute(ctx context.Context, id ScheduledExpenseId) (*ScheduledExpense, error) {
	return action.Repository.FindFor(ctx, id)
}

type FindScheduledExpenses struct {
	Repository ScheduledExpenseRepository
}

func (action *FindScheduledExpenses) Execute(ctx context.Context) ([]ScheduledExpense, error) {
	scheduledExpenses, err := action.Repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(scheduledExpenses, func(i, j int) bool {
		return scheduledExpenses[i].Description < scheduledExpenses[j].Description
	})
	return scheduledExpenses, nil
}

// ErrScheduledExpenseNotFound is returned by Update/Delete/Pause/ResumeScheduledExpense
// when the id doesn't exist in the current user's partition — which, by construction,
// also covers an id owned by someone else. The web layer maps it to 404.
var ErrScheduledExpenseNotFound = errors.New("scheduled expense not found")

type UpdateScheduledExpense struct {
	Repository ScheduledExpenseRepository
}

// Execute preserves Status and LastEvaluatedDate from the existing record —
// neither travels on the wire representation (Status is a list-row action
// only, per ADR 0005; LastEvaluatedDate is generation-engine-internal), and
// Save replaces the whole item, so skipping this would silently reset both
// to their zero values on every edit.
//
// Ownership is enforced structurally, not by comparing UserName: FindFor
// only ever looks inside the current user's own DynamoDB partition
// (PK = user_name from ctx), so an id belonging to another user simply isn't
// found — there is nothing further to compare, unlike revenue/expense whose
// derived composite keys don't encode the owner in the partition key itself.
// See budget-api/CLAUDE.md's Scheduled Expense DynamoDB key scheme.
func (action *UpdateScheduledExpense) Execute(ctx context.Context, scheduledExpense *ScheduledExpense) error {
	applyDefaultTagIfMissing(scheduledExpense)

	existing, err := action.Repository.FindFor(ctx, scheduledExpense.Id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrScheduledExpenseNotFound
	}

	scheduledExpense.Status = existing.Status
	scheduledExpense.LastEvaluatedDate = existing.LastEvaluatedDate
	return action.Repository.Save(ctx, scheduledExpense)
}

type DeleteScheduledExpense struct {
	Repository ScheduledExpenseRepository
}

// Execute hard-deletes only the definition row. BudgetExpenses it already
// generated are left untouched: they carry no reference back to their
// definition (ADR 0005 — traceability lives in Notes only), so there is
// nothing to cascade to.
//
// Ownership is checked at the domain level via FindFor (current user's
// partition only, as in UpdateScheduledExpense); the adapter's
// attribute_exists(id) condition is the backstop for a row that disappears
// between the two calls.
func (action *DeleteScheduledExpense) Execute(ctx context.Context, id ScheduledExpenseId) error {
	existing, err := action.Repository.FindFor(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrScheduledExpenseNotFound
	}
	return action.Repository.Delete(ctx, id)
}

type PauseScheduledExpense struct {
	Repository ScheduledExpenseRepository
	Today      func() date.Date
}

func (action *PauseScheduledExpense) Execute(ctx context.Context, id ScheduledExpenseId) error {
	return transitionStatus(ctx, action.Repository, action.Today, id, StatusPaused)
}

type ResumeScheduledExpense struct {
	Repository ScheduledExpenseRepository
	Today      func() date.Date
}

func (action *ResumeScheduledExpense) Execute(ctx context.Context, id ScheduledExpenseId) error {
	return transitionStatus(ctx, action.Repository, action.Today, id, StatusActive)
}

// transitionStatus moves a definition to target, stamping LastEvaluatedDate
// with today on both pause and resume so no generation gap accumulates across
// a paused span and resuming triggers no backfill (ADR 0005).
//
// Requesting the state a definition is already in is a no-op — no write and,
// crucially, no stamp: re-stamping an already-ACTIVE definition (say, from a
// stale second tab) would silently wipe a pending downtime backfill.
//
// Ownership is checked via FindFor (current user's partition only), as for
// Update/Delete; the adapter's attribute_exists(id) condition is the backstop.
func transitionStatus(ctx context.Context, repository ScheduledExpenseRepository, today func() date.Date, id ScheduledExpenseId, target Status) error {
	existing, err := repository.FindFor(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrScheduledExpenseNotFound
	}
	if existing.Status == target {
		return nil
	}
	return repository.UpdateStatus(ctx, id, target, today())
}
