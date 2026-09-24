package scheduledexpense

import (
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

var aFixedToday = testutils.SafeDateFor("24/09/2026")

func fixedToday() date.Date { return aFixedToday }

func TestWhenAnActiveScheduledExpenseIsPaused(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := PauseScheduledExpense{Repository: mockedRepository, Today: fixedToday}

	existing := &ScheduledExpense{Id: "AN_ID", UserName: "A_USER_NAME", Status: StatusActive}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(existing, nil)
	mockedRepository.On("UpdateStatus", ctx, ScheduledExpenseId("AN_ID"), StatusPaused, aFixedToday).Return(nil)

	err := uut.Execute(ctx, "AN_ID")

	assert.Equal(t, nil, err)
	mockedRepository.AssertExpectations(t)
}

func TestWhenAPausedScheduledExpenseIsResumed(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := ResumeScheduledExpense{Repository: mockedRepository, Today: fixedToday}

	existing := &ScheduledExpense{Id: "AN_ID", UserName: "A_USER_NAME", Status: StatusPaused}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(existing, nil)
	mockedRepository.On("UpdateStatus", ctx, ScheduledExpenseId("AN_ID"), StatusActive, aFixedToday).Return(nil)

	err := uut.Execute(ctx, "AN_ID")

	assert.Equal(t, nil, err)
	mockedRepository.AssertExpectations(t)
}

// Requesting the state a definition is already in is a no-op: no write, and
// in particular no LastEvaluatedDate stamp — re-stamping an ACTIVE definition
// (e.g. from a stale second tab) would silently wipe a pending downtime
// backfill (ADR 0005).
func TestWhenResumingAnAlreadyActiveScheduledExpenseNothingIsWritten(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := ResumeScheduledExpense{Repository: mockedRepository, Today: fixedToday}

	existing := &ScheduledExpense{Id: "AN_ID", UserName: "A_USER_NAME", Status: StatusActive}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(existing, nil)

	err := uut.Execute(ctx, "AN_ID")

	assert.Equal(t, nil, err)
	mockedRepository.AssertNotCalled(t, "UpdateStatus")
}

func TestWhenPausingAnAlreadyPausedScheduledExpenseNothingIsWritten(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := PauseScheduledExpense{Repository: mockedRepository, Today: fixedToday}

	existing := &ScheduledExpense{Id: "AN_ID", UserName: "A_USER_NAME", Status: StatusPaused}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(existing, nil)

	err := uut.Execute(ctx, "AN_ID")

	assert.Equal(t, nil, err)
	mockedRepository.AssertNotCalled(t, "UpdateStatus")
}

// An id not in the caller's own partition — missing, or owned by someone else
// (indistinguishable by construction) — is rejected before any write.
func TestWhenPausingAScheduledExpenseThatDoesNotExist(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := PauseScheduledExpense{Repository: mockedRepository, Today: fixedToday}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("MISSING_ID")).Return(nil, nil)

	err := uut.Execute(ctx, "MISSING_ID")

	assert.Equal(t, ErrScheduledExpenseNotFound, err)
	mockedRepository.AssertNotCalled(t, "UpdateStatus")
}

func TestWhenResumingAScheduledExpenseThatDoesNotExist(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := ResumeScheduledExpense{Repository: mockedRepository, Today: fixedToday}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("MISSING_ID")).Return(nil, nil)

	err := uut.Execute(ctx, "MISSING_ID")

	assert.Equal(t, ErrScheduledExpenseNotFound, err)
	mockedRepository.AssertNotCalled(t, "UpdateStatus")
}

func TestWhenPausingAScheduledExpenseFindForFails(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := PauseScheduledExpense{Repository: mockedRepository, Today: fixedToday}

	ctx := testutils.NewUserContext()
	findError := errors.New("ScheduledExpense FindFor operation fails")
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(nil, findError)

	err := uut.Execute(ctx, "AN_ID")

	assert.Equal(t, findError, err)
	mockedRepository.AssertNotCalled(t, "UpdateStatus")
}

func TestWhenResumingAScheduledExpenseUpdateStatusFails(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := ResumeScheduledExpense{Repository: mockedRepository, Today: fixedToday}

	existing := &ScheduledExpense{Id: "AN_ID", UserName: "A_USER_NAME", Status: StatusPaused}

	ctx := testutils.NewUserContext()
	updateError := errors.New("ScheduledExpense UpdateStatus operation fails")
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(existing, nil)
	mockedRepository.On("UpdateStatus", ctx, ScheduledExpenseId("AN_ID"), StatusActive, aFixedToday).Return(updateError)

	err := uut.Execute(ctx, "AN_ID")

	assert.Equal(t, updateError, err)
}
