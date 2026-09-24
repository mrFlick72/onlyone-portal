package scheduledexpense

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

var generationTestLogger = logging.GetLoggerInstanceForComponentByTypeName("scheduledexpense-generation-test")

func iso(s string) date.Date {
	d, err := date.IsoDateFor(s)
	if err != nil {
		panic(err)
	}
	return *d
}

func isoPtr(s string) *date.Date {
	d := iso(s)
	return &d
}

func intPtr(i int) *int { return &i }

// generationRecorder logs creates and advances, in order, across both fakes —
// so tests can assert generate-then-advance ordering.
type generationRecorder struct {
	calls []string
}

func (r *generationRecorder) record(format string, args ...any) {
	r.calls = append(r.calls, fmt.Sprintf(format, args...))
}

type fakeGenerationRepository struct {
	rec         *generationRecorder
	definitions []ScheduledExpense
	findErr     error
	// advanceErrFor makes AdvanceLastEvaluatedDate fail for the given ISO day
	// (once), simulating a crash between generate and advance.
	advanceErrFor map[string]error
	advanced      map[ScheduledExpenseId]date.Date
}

func (f *fakeGenerationRepository) FindAllActive(_ context.Context) ([]ScheduledExpense, error) {
	return f.definitions, f.findErr
}

func (f *fakeGenerationRepository) AdvanceLastEvaluatedDate(ctx context.Context, id ScheduledExpenseId, d date.Date) error {
	user, _ := security.GetCurrentUser(ctx)
	f.rec.record("advance %s %s %s", *user.UserName, id, d.GetIsoFormattedDate())
	if err, ok := f.advanceErrFor[d.GetIsoFormattedDate()]; ok {
		delete(f.advanceErrFor, d.GetIsoFormattedDate())
		return err
	}
	if f.advanced == nil {
		f.advanced = map[ScheduledExpenseId]date.Date{}
	}
	f.advanced[id] = d
	// Mirror the store: the next run sees the advanced date.
	for i := range f.definitions {
		if f.definitions[i].Id == id {
			advanced := d
			f.definitions[i].LastEvaluatedDate = &advanced
		}
	}
	return nil
}

type fakeExpenseCreator struct {
	rec       *generationRecorder
	created   []expense.BudgetExpense
	users     []security.User
	createErr error
	panicFor  string
}

func (f *fakeExpenseCreator) CreateBudgetExpense(ctx context.Context, budgetExpense *expense.BudgetExpense) error {
	user, _ := security.GetCurrentUser(ctx)
	f.rec.record("create %s %s", *user.UserName, budgetExpense.Date.GetIsoFormattedDate())
	if f.panicFor != "" && *user.UserName == f.panicFor {
		panic("boom")
	}
	if f.createErr != nil {
		return f.createErr
	}
	f.created = append(f.created, *budgetExpense)
	f.users = append(f.users, *user)
	return nil
}

func newGeneration(today string, definitions ...ScheduledExpense) (*GenerateScheduledExpenses, *fakeGenerationRepository, *fakeExpenseCreator, *generationRecorder) {
	rec := &generationRecorder{}
	repository := &fakeGenerationRepository{rec: rec, definitions: definitions, advanceErrFor: map[string]error{}}
	creator := &fakeExpenseCreator{rec: rec}
	uut := &GenerateScheduledExpenses{
		Repository:     repository,
		ExpenseCreator: creator,
		Today:          func() date.Date { return iso(today) },
		Logger:         generationTestLogger,
	}
	return uut, repository, creator, rec
}

func createdDates(creator *fakeExpenseCreator) []string {
	dates := []string{}
	for _, e := range creator.created {
		dates = append(dates, e.Date.GetIsoFormattedDate())
	}
	return dates
}

func activeDefinition(day int) ScheduledExpense {
	return ScheduledExpense{
		Id:          "SE_ID",
		UserName:    "owner",
		Description: "Rent",
		Amount:      testutils.SafeMoneyFor("1200.00"),
		Tags:        []tags.SearchTag{{Key: "housing", Value: "Housing"}},
		Day:         day,
		Status:      StatusActive,
	}
}

// Day/Month matching, month-end clamping, End Date and backfill are all a
// function of (definition, LastEvaluatedDate, today) → generated days.
func TestGenerationGeneratesExactlyTheMatchingDays(t *testing.T) {
	cases := []struct {
		name          string
		day           int
		month         *int
		endDate       *date.Date
		lastEvaluated *date.Date
		today         string
		expected      []string
	}{
		{name: "exact day match, monthly", day: 5, lastEvaluated: isoPtr("2026-09-04"), today: "2026-09-05", expected: []string{"2026-09-05"}},
		{name: "not the day", day: 6, lastEvaluated: isoPtr("2026-09-04"), today: "2026-09-05", expected: []string{}},
		{name: "day 31 clamps to 30 April", day: 31, lastEvaluated: isoPtr("2026-04-29"), today: "2026-04-30", expected: []string{"2026-04-30"}},
		{name: "day 31 does not fire on 29 April", day: 31, lastEvaluated: isoPtr("2026-04-28"), today: "2026-04-29", expected: []string{}},
		{name: "day 29 clamps to 28 Feb (non-leap)", day: 29, lastEvaluated: isoPtr("2026-02-27"), today: "2026-02-28", expected: []string{"2026-02-28"}},
		{name: "day 31 clamps to 28 Feb (non-leap)", day: 31, lastEvaluated: isoPtr("2026-02-27"), today: "2026-02-28", expected: []string{"2026-02-28"}},
		{name: "day 31 clamps to 29 Feb (leap), not the 28th", day: 31, lastEvaluated: isoPtr("2028-02-27"), today: "2028-02-29", expected: []string{"2028-02-29"}},
		{name: "day 29 fires on 29 Feb in a leap year", day: 29, lastEvaluated: isoPtr("2028-02-28"), today: "2028-02-29", expected: []string{"2028-02-29"}},
		{name: "yearly fires on its month", day: 15, month: intPtr(3), lastEvaluated: isoPtr("2026-03-14"), today: "2026-03-15", expected: []string{"2026-03-15"}},
		{name: "yearly skips other months", day: 15, month: intPtr(3), lastEvaluated: isoPtr("2026-04-14"), today: "2026-04-15", expected: []string{}},
		{name: "monthly fires every month", day: 15, lastEvaluated: isoPtr("2026-04-14"), today: "2026-04-15", expected: []string{"2026-04-15"}},
		{name: "End Date day itself still fires", day: 5, endDate: isoPtr("2026-09-05"), lastEvaluated: isoPtr("2026-09-04"), today: "2026-09-05", expected: []string{"2026-09-05"}},
		{name: "past End Date never fires", day: 5, endDate: isoPtr("2026-09-04"), lastEvaluated: isoPtr("2026-09-04"), today: "2026-09-05", expected: []string{}},
		// Resume stamps LastEvaluatedDate with the resume day (ADR 0005), so a
		// definition resumed after its End Date is simply inert — no special case.
		{name: "resumed after End Date stays inert", day: 5, endDate: isoPtr("2026-06-30"), lastEvaluated: isoPtr("2026-10-04"), today: "2026-10-05", expected: []string{}},
		{name: "backfill across a multi-month gap", day: 5, lastEvaluated: isoPtr("2026-07-01"), today: "2026-09-24", expected: []string{"2026-07-05", "2026-08-05", "2026-09-05"}},
		{name: "backfill clamps per month", day: 31, lastEvaluated: isoPtr("2026-08-01"), today: "2026-09-30", expected: []string{"2026-08-31", "2026-09-30"}},
		{name: "backfill capped at End Date", day: 5, endDate: isoPtr("2026-08-10"), lastEvaluated: isoPtr("2026-07-01"), today: "2026-09-24", expected: []string{"2026-07-05", "2026-08-05"}},
		{name: "never evaluated: today only, matching", day: 24, today: "2026-09-24", expected: []string{"2026-09-24"}},
		{name: "never evaluated: today only, no backfill", day: 5, today: "2026-09-24", expected: []string{}},
		{name: "already evaluated today", day: 24, lastEvaluated: isoPtr("2026-09-24"), today: "2026-09-24", expected: []string{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			definition := activeDefinition(c.day)
			definition.Month = c.month
			definition.EndDate = c.endDate
			definition.LastEvaluatedDate = c.lastEvaluated
			uut, repository, creator, _ := newGeneration(c.today, definition)

			err := uut.Execute(context.Background())

			assert.Equal(t, nil, err)
			assert.Equal(t, c.expected, createdDates(creator))
			if c.lastEvaluated != nil && c.lastEvaluated.GetIsoFormattedDate() == c.today {
				assert.Equal(t, 0, len(repository.advanced))
				return
			}
			advanced := repository.advanced["SE_ID"]
			assert.Equal(t, c.today, advanced.GetIsoFormattedDate())
		})
	}
}

// Each generated day is created before LastEvaluatedDate advances to it; the
// walk then advances to today once.
func TestGenerationCreatesBeforeAdvancingEachDay(t *testing.T) {
	definition := activeDefinition(5)
	definition.LastEvaluatedDate = isoPtr("2026-08-01")
	uut, _, _, rec := newGeneration("2026-09-24", definition)

	err := uut.Execute(context.Background())

	assert.Equal(t, nil, err)
	assert.Equal(t, []string{
		"create owner 2026-08-05",
		"advance owner SE_ID 2026-08-05",
		"create owner 2026-09-05",
		"advance owner SE_ID 2026-09-05",
		"advance owner SE_ID 2026-09-24",
	}, rec.calls)
}

// A crash between generate and advance must surface as a visible duplicate on
// the next run, never as a silently lost expense (ADR 0005).
func TestGenerationCrashBetweenGenerateAndAdvanceDuplicatesRatherThanLoses(t *testing.T) {
	definition := activeDefinition(5)
	definition.LastEvaluatedDate = isoPtr("2026-09-04")
	uut, repository, creator, _ := newGeneration("2026-09-05", definition)
	repository.advanceErrFor["2026-09-05"] = errors.New("crash before advance")

	_ = uut.Execute(context.Background())
	assert.Equal(t, []string{"2026-09-05"}, createdDates(creator))
	assert.Equal(t, "2026-09-04", repository.definitions[0].LastEvaluatedDate.GetIsoFormattedDate())

	_ = uut.Execute(context.Background())

	assert.Equal(t, []string{"2026-09-05", "2026-09-05"}, createdDates(creator))
	assert.Equal(t, "2026-09-05", repository.definitions[0].LastEvaluatedDate.GetIsoFormattedDate())
}

// A failed create stops that definition's walk without advancing past the
// failed day (so the next run retries it), and doesn't stop the others.
func TestGenerationCreateFailureDoesNotAdvanceAndOtherDefinitionsStillRun(t *testing.T) {
	failing := activeDefinition(5)
	failing.LastEvaluatedDate = isoPtr("2026-09-04")
	other := activeDefinition(5)
	other.Id = "OTHER_ID"
	other.UserName = "other"
	other.LastEvaluatedDate = isoPtr("2026-09-04")
	uut, repository, creator, rec := newGeneration("2026-09-05", failing, other)
	creator.createErr = errors.New("DynamoDB down")

	err := uut.Execute(context.Background())

	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"create owner 2026-09-05", "create other 2026-09-05"}, rec.calls)
	assert.Equal(t, 0, len(repository.advanced))
}

// A panic while generating one definition is recovered, so it can't take the
// whole of budget-api down from a background goroutine; the rest still run.
func TestGenerationRecoversFromAPanicInOneDefinition(t *testing.T) {
	panicking := activeDefinition(5)
	panicking.UserName = "panics"
	panicking.LastEvaluatedDate = isoPtr("2026-09-04")
	other := activeDefinition(5)
	other.Id = "OTHER_ID"
	other.UserName = "other"
	other.LastEvaluatedDate = isoPtr("2026-09-04")
	uut, _, creator, _ := newGeneration("2026-09-05", panicking, other)
	creator.panicFor = "panics"

	err := uut.Execute(context.Background())

	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"2026-09-05"}, createdDates(creator))
	assert.Equal(t, "other", *creator.users[0].UserName)
}

// Paused definitions are skipped entirely — not evaluated, not advanced —
// even if one reaches the engine (seeded directly, not via PATCH).
func TestGenerationSkipsPausedDefinitions(t *testing.T) {
	paused := activeDefinition(5)
	paused.Status = StatusPaused
	paused.LastEvaluatedDate = isoPtr("2026-08-01")
	uut, repository, creator, rec := newGeneration("2026-09-24", paused)

	err := uut.Execute(context.Background())

	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(creator.created))
	assert.Equal(t, 0, len(repository.advanced))
	assert.Equal(t, 0, len(rec.calls))
}

// The generated expense: the matching day's date, the definition's amount and
// tags (keys and stored names), and Note = Notes + a trace line naming the
// definition. Created under an owner-only context (no token/authorities).
func TestGeneratedExpenseContent(t *testing.T) {
	definition := activeDefinition(24)
	definition.Id = "1a2b3c"
	definition.Notes = "Monthly rent, paid by bank transfer"
	uut, _, creator, _ := newGeneration("2026-09-24", definition)

	err := uut.Execute(context.Background())

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(creator.created))
	generated := creator.created[0]
	assert.Equal(t, "", generated.Id)
	assert.Equal(t, "2026-09-24", generated.Date.GetIsoFormattedDate())
	assert.Equal(t, "1200.00", generated.Amount.StringifyAmount())
	assert.Equal(t, []tags.SearchTag{{Key: "housing", Value: "Housing"}}, generated.Tags)
	assert.Equal(t, "Monthly rent, paid by bank transfer\nExpense generated by the scheduled expense \"Rent\" with id: 1a2b3c", generated.Note)
	assert.Equal(t, "owner", *creator.users[0].UserName)
	assert.Equal(t, true, creator.users[0].AccessToken == nil)
	assert.Equal(t, true, creator.users[0].Authorities == nil)
}

func TestGeneratedExpenseNoteWithoutNotesIsJustTheTraceLine(t *testing.T) {
	definition := activeDefinition(24)
	definition.Id = "1a2b3c"
	uut, _, creator, _ := newGeneration("2026-09-24", definition)

	_ = uut.Execute(context.Background())

	assert.Equal(t, "Expense generated by the scheduled expense \"Rent\" with id: 1a2b3c", creator.created[0].Note)
}

// Shutdown cancels the job's context: the walk stops between days.
func TestGenerationStopsWhenTheContextIsCancelled(t *testing.T) {
	definition := activeDefinition(5)
	definition.LastEvaluatedDate = isoPtr("2026-07-01")
	uut, _, creator, _ := newGeneration("2026-09-24", definition)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := uut.Execute(ctx)

	assert.Equal(t, context.Canceled, err)
	assert.Equal(t, 0, len(creator.created))
}

func TestGenerationReturnsTheErrorWhenListingDefinitionsFails(t *testing.T) {
	uut, repository, _, _ := newGeneration("2026-09-24")
	findErr := errors.New("scan fails")
	repository.findErr = findErr

	err := uut.Execute(context.Background())

	assert.Equal(t, findErr, err)
}
