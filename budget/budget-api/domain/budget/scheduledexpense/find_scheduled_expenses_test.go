package scheduledexpense

import (
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestWhenScheduledExpensesAreFoundTheyAreSortedByDescription(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := FindScheduledExpenses{Repository: mockedRepository}

	netflix := ScheduledExpense{Id: "2", Description: "Netflix subscription", Day: 1}
	insurance := ScheduledExpense{Id: "1", Description: "Insurance", Day: 15}
	rent := ScheduledExpense{Id: "3", Description: "Rent", Day: 5}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindAll", ctx).Return([]ScheduledExpense{netflix, insurance, rent}, nil)

	result, err := uut.Execute(ctx)

	assert.Equal(t, nil, err)
	assert.Equal(t, []ScheduledExpense{insurance, netflix, rent}, result)
}

func TestWhenNoScheduledExpensesExistAnEmptyListIsReturned(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := FindScheduledExpenses{Repository: mockedRepository}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindAll", ctx).Return([]ScheduledExpense{}, nil)

	result, err := uut.Execute(ctx)

	assert.Equal(t, nil, err)
	assert.Equal(t, []ScheduledExpense{}, result)
}

func TestWhenFindingScheduledExpensesFails(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := FindScheduledExpenses{Repository: mockedRepository}

	ctx := testutils.NewUserContext()
	findError := errors.New("ScheduledExpense FindAll operation fails")
	mockedRepository.On("FindAll", ctx).Return(nil, findError)

	result, err := uut.Execute(ctx)

	assert.Equal(t, findError, err)
	assert.Equal(t, 0, len(result))
}
