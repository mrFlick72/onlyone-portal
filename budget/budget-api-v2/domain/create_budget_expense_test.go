package domain

import (
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestWhenANewBudgetExpenseISCreated(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	uut := CreateBudgetExpense{
		repository: mockedRepository,
	}

	aDate, _ := IsoDateFor("2018-01-01")
	anAmount, _ := MoneyFor("1.00")
	aBudgetExpense := BudgetExpense{
		Date:   *aDate,
		Amount: *anAmount,
		Note:   "A_NOTE",
		Tag:    "super-market",
	}

	ctx := newUserContext()

	mockedRepository.On("Save", &ctx, &aBudgetExpense).Return(nil)

	err := uut.Execute(&ctx, &aBudgetExpense)

	assert.Equal(t, nil, err)
	assert.Equal(t, "A_USER_NAME", aBudgetExpense.UserName)
	assert.Equal(t, "A_BUDGET_ID", aBudgetExpense.Id)
	mockedRepository.AssertCalled(t, "Save", &ctx, &aBudgetExpense)
}

func TestWhenANewBudgetExpenseCreationFails(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	uut := CreateBudgetExpense{
		repository: mockedRepository,
	}

	aDate, _ := IsoDateFor("2018-01-01")
	anAmount, _ := MoneyFor("1.00")
	aBudgetExpense := BudgetExpense{
		Date:   *aDate,
		Amount: *anAmount,
		Note:   "A_NOTE",
		Tag:    "super-market",
	}

	ctx := newUserContext()

	saveError := errors.New("Budget Expense Save operation fails")
	mockedRepository.On("Save", &ctx, &aBudgetExpense).Return(saveError)

	err := uut.Execute(&ctx, &aBudgetExpense)

	assert.Equal(t, saveError, err)
	mockedRepository.AssertCalled(t, "Save", &ctx, &aBudgetExpense)
}