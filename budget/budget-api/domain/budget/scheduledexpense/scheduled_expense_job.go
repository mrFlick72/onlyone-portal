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

// ScheduledExpenseJob is the scheduled job behind Scheduled Expenses: each
// Execute evaluates every ACTIVE definition for each day from
// LastEvaluatedDate+1 (today, when never evaluated) up to today, stores a
// BudgetExpense for each matching day and advances the definition's
// LastEvaluatedDate. Running it more than once a day is safe:
// LastEvaluatedDate makes each day evaluated exactly once. It is run on a
// timer by adapter/budget/scheduledexpense/scheduler's
// ScheduledExpenseJobConfigurer. See ADR 0005.
//
// Expenses are created through the existing CreateBudgetExpense action — the
// same instance the expense endpoints use — so a generated expense goes
// through exactly the path a manual one does (default tag, Kafka event).
type ScheduledExpenseJob struct {
	Repository          ScheduledExpenseRepository
	CreateBudgetExpense *expense.CreateBudgetExpense
	Today               func() date.Date
	Logger              *logging.Logger
}

// Execute returns an error only when the definitions can't be listed or ctx
// is cancelled; a failure on one definition is logged and the run moves on.
func (job *ScheduledExpenseJob) Execute(ctx context.Context) error {
	definitions, err := job.Repository.FindAllActive(ctx)
	if err != nil {
		job.Logger.LogErrorfFor("scheduled expense job: listing active definitions failed: %v", err)
		return err
	}

	today := job.Today()
	for _, definition := range definitions {
		if err := ctx.Err(); err != nil {
			return err
		}
		if definition.Status != StatusActive {
			continue
		}
		if err := job.generateFor(ctx, definition, today); err != nil && ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return nil
}

// generateFor walks one definition's pending days. It recovers from a panic
// so a single bad definition can't take budget-api down from the scheduler's
// goroutine.
func (job *ScheduledExpenseJob) generateFor(ctx context.Context, definition ScheduledExpense, today date.Date) (err error) {
	defer func() {
		if r := recover(); r != nil {
			job.Logger.LogErrorfFor("scheduled expense job: panic on definition %s: %v", definition.Id, r)
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
		if err := job.CreateBudgetExpense.Execute(ownerCtx, definition.expenseFor(day)); err != nil {
			job.Logger.LogErrorfFor("scheduled expense job: creating expense for definition %s on %s failed: %v", definition.Id, day.GetIsoFormattedDate(), err)
			return err
		}
		if err := job.advance(ownerCtx, definition, day); err != nil {
			return err
		}
		advanced := day
		lastAdvanced = &advanced
	}

	if lastAdvanced == nil || today.IsAfter(*lastAdvanced) {
		return job.advance(ownerCtx, definition, today)
	}
	return nil
}

func (job *ScheduledExpenseJob) advance(ownerCtx context.Context, definition ScheduledExpense, day date.Date) error {
	err := job.Repository.AdvanceLastEvaluatedDate(ownerCtx, definition.Id, day)
	if errors.Is(err, ErrScheduledExpenseNotFound) {
		job.Logger.LogInfofFor("scheduled expense job: definition %s was deleted or paused mid-run, stopping", definition.Id)
		return err
	}
	if err != nil {
		job.Logger.LogErrorfFor("scheduled expense job: advancing definition %s to %s failed: %v", definition.Id, day.GetIsoFormattedDate(), err)
	}
	return err
}
