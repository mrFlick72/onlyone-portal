package scheduledexpense

import (
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestWhenAScheduledExpenseIsFoundById(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := FindScheduledExpense{Repository: mockedRepository}

	expected := &ScheduledExpense{Id: "AN_ID", Description: "Rent", Day: 5}
	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(expected, nil)

	result, err := uut.Execute(ctx, "AN_ID")

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, result)
}

func TestWhenAScheduledExpenseIsNotFoundById(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := FindScheduledExpense{Repository: mockedRepository}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("MISSING_ID")).Return(nil, nil)

	result, err := uut.Execute(ctx, "MISSING_ID")

	assert.Equal(t, nil, err)
	assert.Equal(t, true, result == nil)
}

func TestWhenFindingAScheduledExpenseByIdFails(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := FindScheduledExpense{Repository: mockedRepository}

	ctx := testutils.NewUserContext()
	findError := errors.New("ScheduledExpense FindFor operation fails")
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(nil, findError)

	result, err := uut.Execute(ctx, "AN_ID")

	assert.Equal(t, findError, err)
	assert.Equal(t, true, result == nil)
}
