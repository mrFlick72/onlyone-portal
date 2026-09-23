package scheduledexpense

import (
	"context"
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
