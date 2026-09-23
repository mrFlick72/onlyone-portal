package scheduledexpense

import (
	"context"
	"errors"
	"sort"

	"github.com/mrflick72/budget/budget-api/domain/tags"
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
		return errors.New("scheduled expense not found or user not authorized to update it")
	}

	scheduledExpense.Status = existing.Status
	scheduledExpense.LastEvaluatedDate = existing.LastEvaluatedDate
	return action.Repository.Save(ctx, scheduledExpense)
}
