package scheduledexpense

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestWhenANewScheduledExpenseIsCreated(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := CreateScheduledExpense{Repository: mockedRepository}

	anAmount := testutils.SafeMoneyFor("100.00")
	aScheduledExpense := ScheduledExpense{
		Description: "Rent",
		Amount:      anAmount,
		Notes:       "A_NOTE",
		Day:         5,
	}

	ctx := testutils.NewUserContext()
	mockedRepository.On("Save", ctx, &aScheduledExpense).Return(nil)

	err := uut.Execute(ctx, &aScheduledExpense)

	assert.Equal(t, nil, err)
	assert.Equal(t, "A_USER_NAME", aScheduledExpense.UserName)
	assert.Equal(t, "A_SCHEDULED_EXPENSE_ID", aScheduledExpense.Id)
	assert.Equal(t, StatusActive, aScheduledExpense.Status)
	mockedRepository.AssertCalled(t, "Save", ctx, &aScheduledExpense)
}

func TestWhenANewScheduledExpenseHasNoTagsItDefaultsToUnknown(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := CreateScheduledExpense{Repository: mockedRepository}

	anAmount := testutils.SafeMoneyFor("100.00")
	aScheduledExpense := ScheduledExpense{Description: "Rent", Amount: anAmount, Notes: "A_NOTE", Day: 5}

	ctx := testutils.NewUserContext()
	mockedRepository.On("Save", ctx, &aScheduledExpense).Return(nil)

	err := uut.Execute(ctx, &aScheduledExpense)

	assert.Equal(t, nil, err)
	assert.Equal(t, []tags.SearchTag{tags.UnknownSentinel()}, aScheduledExpense.Tags)
}

func TestWhenANewScheduledExpenseHasTagsTheyArePreserved(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := CreateScheduledExpense{Repository: mockedRepository}

	anAmount := testutils.SafeMoneyFor("100.00")
	provided := []tags.SearchTag{{Key: "housing", Value: "Housing"}}
	aScheduledExpense := ScheduledExpense{Description: "Rent", Amount: anAmount, Notes: "A_NOTE", Day: 5, Tags: provided}

	ctx := testutils.NewUserContext()
	mockedRepository.On("Save", ctx, &aScheduledExpense).Return(nil)

	err := uut.Execute(ctx, &aScheduledExpense)

	assert.Equal(t, nil, err)
	assert.Equal(t, provided, aScheduledExpense.Tags)
}

func TestWhenANewScheduledExpenseCreationFailsBecauseNoUserInContext(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := CreateScheduledExpense{Repository: mockedRepository}

	anAmount := testutils.SafeMoneyFor("100.00")
	aScheduledExpense := ScheduledExpense{Description: "Rent", Amount: anAmount, Notes: "A_NOTE", Day: 5}

	err := uut.Execute(context.Background(), &aScheduledExpense)

	assert.NotEqual(t, nil, err)
	mockedRepository.AssertNotCalled(t, "Save")
}

func TestWhenANewScheduledExpenseCreationFails(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := CreateScheduledExpense{Repository: mockedRepository}

	anAmount := testutils.SafeMoneyFor("100.00")
	aScheduledExpense := ScheduledExpense{Description: "Rent", Amount: anAmount, Notes: "A_NOTE", Day: 5}

	ctx := testutils.NewUserContext()
	saveError := errors.New("ScheduledExpense Save operation fails")
	mockedRepository.On("Save", ctx, &aScheduledExpense).Return(saveError)

	err := uut.Execute(ctx, &aScheduledExpense)

	assert.Equal(t, saveError, err)
	mockedRepository.AssertCalled(t, "Save", ctx, &aScheduledExpense)
}
