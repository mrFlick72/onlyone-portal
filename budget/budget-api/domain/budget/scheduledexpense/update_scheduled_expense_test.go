package scheduledexpense

import (
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestWhenAnExistingScheduledExpenseIsUpdated(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := UpdateScheduledExpense{Repository: mockedRepository}

	lastEvaluated := testutils.SafeDateFor("01/09/2026")
	existing := &ScheduledExpense{
		Id:                "AN_ID",
		UserName:          "A_USER_NAME",
		Description:       "Old description",
		Day:               1,
		Status:            StatusPaused,
		LastEvaluatedDate: &lastEvaluated,
	}

	update := ScheduledExpense{
		Id:          "AN_ID",
		Description: "New description",
		Amount:      testutils.SafeMoneyFor("50.00"),
		Day:         10,
	}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(existing, nil)
	mockedRepository.On("Save", ctx, &update).Return(nil)

	err := uut.Execute(ctx, &update)

	assert.Equal(t, nil, err)
	// Status and LastEvaluatedDate are carried over from the existing record —
	// they don't travel on the wire representation, and Save replaces the
	// whole item, so skipping this would silently reset both.
	assert.Equal(t, StatusPaused, update.Status)
	assert.Equal(t, &lastEvaluated, update.LastEvaluatedDate)
	mockedRepository.AssertCalled(t, "Save", ctx, &update)
}

func TestWhenAnUpdatedScheduledExpenseHasNoTagsItDefaultsToUnknown(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := UpdateScheduledExpense{Repository: mockedRepository}

	existing := &ScheduledExpense{Id: "AN_ID", UserName: "A_USER_NAME", Status: StatusActive}
	update := ScheduledExpense{Id: "AN_ID", Description: "Rent", Day: 5}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(existing, nil)
	mockedRepository.On("Save", ctx, &update).Return(nil)

	err := uut.Execute(ctx, &update)

	assert.Equal(t, nil, err)
	assert.Equal(t, []tags.SearchTag{tags.UnknownSentinel()}, update.Tags)
}

func TestWhenUpdatingAScheduledExpenseThatDoesNotExist(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := UpdateScheduledExpense{Repository: mockedRepository}

	update := ScheduledExpense{Id: "MISSING_ID", Description: "Rent", Day: 5}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("MISSING_ID")).Return(nil, nil)

	err := uut.Execute(ctx, &update)

	assert.Equal(t, ErrScheduledExpenseNotFound, err)
	mockedRepository.AssertNotCalled(t, "Save")
}

func TestWhenUpdatingAScheduledExpenseFindForFails(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := UpdateScheduledExpense{Repository: mockedRepository}

	update := ScheduledExpense{Id: "AN_ID", Description: "Rent", Day: 5}

	ctx := testutils.NewUserContext()
	findError := errors.New("ScheduledExpense FindFor operation fails")
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(nil, findError)

	err := uut.Execute(ctx, &update)

	assert.Equal(t, findError, err)
	mockedRepository.AssertNotCalled(t, "Save")
}

func TestWhenUpdatingAScheduledExpenseSaveFails(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := UpdateScheduledExpense{Repository: mockedRepository}

	existing := &ScheduledExpense{Id: "AN_ID", UserName: "A_USER_NAME", Status: StatusActive}
	update := ScheduledExpense{Id: "AN_ID", Description: "Rent", Day: 5}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(existing, nil)
	saveError := errors.New("ScheduledExpense Save operation fails")
	mockedRepository.On("Save", ctx, &update).Return(saveError)

	err := uut.Execute(ctx, &update)

	assert.Equal(t, saveError, err)
}

