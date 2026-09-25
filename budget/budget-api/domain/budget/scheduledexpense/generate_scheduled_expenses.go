package scheduledexpense

import (
	"context"
	"errors"
	"fmt"

	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

// GenerateScheduledExpenses is one run of the generation engine: it evaluates
// every ACTIVE definition for each day from LastEvaluatedDate+1 (today, when
// never evaluated) up to today, and creates a BudgetExpense for each matching
// day. Running it more than once a day is safe: LastEvaluatedDate makes each
// day evaluated exactly once. See ADR 0005.
//
// Expenses are created through the existing CreateBudgetExpense action — the
// same instance the expense endpoints use — so a generated expense goes
// through exactly the path a manual one does (default tag, Kafka event).
type GenerateScheduledExpenses struct {
	Repository          ScheduledExpenseRepository
	CreateBudgetExpense *expense.CreateBudgetExpense
	Today               func() date.Date
	Logger              *logging.Logger
}

// Execute returns an error only when the definitions can't be listed or ctx
// is cancelled; a failure on one definition is logged and the run moves on.
func (action *GenerateScheduledExpenses) Execute(ctx context.Context) error {
	definitions, err := action.Repository.FindAllActive(ctx)
	if err != nil {
		action.Logger.LogErrorfFor("scheduled expense generation: listing active definitions failed: %v", err)
		return err
	}

	today := action.Today()
	for _, definition := range definitions {
		if err := ctx.Err(); err != nil {
			return err
		}
		if definition.Status != StatusActive {
			continue
		}
		if err := action.generateFor(ctx, definition, today); err != nil && ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return nil
}

// generateFor walks one definition's pending days. It recovers from a panic
// so a single bad definition can't take budget-api down from the scheduler's
// goroutine.
func (action *GenerateScheduledExpenses) generateFor(ctx context.Context, definition ScheduledExpense, today date.Date) (err error) {
	defer func() {
		if r := recover(); r != nil {
			action.Logger.LogErrorfFor("scheduled expense generation: panic on definition %s: %v", definition.Id, r)
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	// Owner-only context: CreateBudgetExpense and the Kafka publisher only
	// read UserName (ADR 0005). No AccessToken, no Authorities.
	owner := definition.UserName
	ownerCtx := context.WithValue(ctx, "user", security.User{UserName: &owner})

	day := today
	if definition.LastEvaluatedDate != nil {
		day = definition.LastEvaluatedDate.AddDays(1)
	}
	if day.IsAfter(today) {
		return nil
	}

	var lastAdvanced *date.Date
	for ; !day.IsAfter(today); day = day.AddDays(1) {
		if err := ctx.Err(); err != nil {
			return err
		}
		if !definition.isDueOn(day) {
			continue
		}
		// Generate, then advance: a crash in between re-evaluates this day on
		// the next run — a visible duplicate, never a silent loss.
		if err := action.CreateBudgetExpense.Execute(ownerCtx, definition.expenseFor(day)); err != nil {
			action.Logger.LogErrorfFor("scheduled expense generation: creating expense for definition %s on %s failed: %v", definition.Id, day.GetIsoFormattedDate(), err)
			return err
		}
		if err := action.advance(ownerCtx, definition, day); err != nil {
			return err
		}
		advanced := day
		lastAdvanced = &advanced
	}

	if lastAdvanced == nil || today.IsAfter(*lastAdvanced) {
		return action.advance(ownerCtx, definition, today)
	}
	return nil
}

func (action *GenerateScheduledExpenses) advance(ownerCtx context.Context, definition ScheduledExpense, day date.Date) error {
	err := action.Repository.AdvanceLastEvaluatedDate(ownerCtx, definition.Id, day)
	if errors.Is(err, ErrScheduledExpenseNotFound) {
		action.Logger.LogInfofFor("scheduled expense generation: definition %s was deleted or paused mid-run, stopping", definition.Id)
		return err
	}
	if err != nil {
		action.Logger.LogErrorfFor("scheduled expense generation: advancing definition %s to %s failed: %v", definition.Id, day.GetIsoFormattedDate(), err)
	}
	return err
}
