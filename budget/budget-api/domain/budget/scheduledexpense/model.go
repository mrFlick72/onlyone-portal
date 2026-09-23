package scheduledexpense

import (
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
)

// ScheduledExpense is a per-user template that the daily generation engine
// evaluates and, when due, turns into a real expense.BudgetExpense for its
// owner. See budget/budget-api/CONTEXT.md ("Scheduled Expense") and
// docs/adr/0005-scheduled-expense-recurrence-and-generation-engine.md for the
// full domain model and the trade-offs behind Day/Month/EndDate/Status.
type ScheduledExpense struct {
	Id       ScheduledExpenseId
	UserName UserName

	// Description is the required, user-facing label identifying the purpose
	// of this template (e.g. "Rent", "Netflix subscription").
	Description string
	Amount      money.Money
	Notes       string
	Tags        []tags.SearchTag

	// Day is the day of month (1-31) the recurrence fires on. Month, when set,
	// makes the recurrence yearly; left nil, it recurs monthly.
	Day   int
	Month *int

	// EndDate bounds the recurrence: nil means it recurs forever, a set date
	// means the daily job stops generating once that date has passed.
	EndDate *date.Date

	Status Status

	// LastEvaluatedDate is stamped by the daily generation engine on every
	// evaluation (including pause/resume), driving its downtime-backfill logic.
	// Nil until the definition has been evaluated for the first time.
	LastEvaluatedDate *date.Date
}

type ScheduledExpenseId = string
type UserName = string

// Status is the pause/resume state of a ScheduledExpense. The zero value ""
// is neither ACTIVE nor PAUSED, so an unset Status is never mistaken for a
// deliberately chosen one.
type Status string

const (
	StatusActive Status = "ACTIVE"
	StatusPaused Status = "PAUSED"
)

type ScheduledExpenseIdProvider interface {
	GenerateIdFor(scheduledExpense *ScheduledExpense) ScheduledExpenseId
}
