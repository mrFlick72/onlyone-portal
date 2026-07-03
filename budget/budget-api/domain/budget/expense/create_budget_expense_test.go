package expense

import (
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/stretchr/testify/mock"
)

func TestWhenANewBudgetExpenseISCreated(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	mockedPublisher := new(BudgetExpenseEventPublisherMock)
	uut := CreateBudgetExpense{
		Repository:     mockedRepository,
		EventPublisher: mockedPublisher,
	}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aBudgetExpense := BudgetExpense{
		Date:   *aDate,
		Amount: anAmount,
		Note:   "A_NOTE",
		Tags:   []tags.SearchTag{{Key: "super-market", Value: "super-market"}},
	}

	ctx := testutils.NewUserContext()

	mockedRepository.On("Save", ctx, &aBudgetExpense).Return(nil)
	mockedPublisher.On("CreateBudgetExpense", ctx, mock.Anything).Return(nil)

	err := uut.Execute(ctx, &aBudgetExpense)

	assert.Equal(t, nil, err)
	assert.Equal(t, "A_USER_NAME", aBudgetExpense.UserName)
	assert.Equal(t, "A_BUDGET_ID", aBudgetExpense.Id)
	mockedRepository.AssertCalled(t, "Save", ctx, &aBudgetExpense)
	mockedPublisher.AssertCalled(t, "CreateBudgetExpense", ctx, aBudgetExpense)
}

func TestWhenANewBudgetExpenseWithNoTagsDefaultsToUnknown(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	mockedPublisher := new(BudgetExpenseEventPublisherMock)
	uut := CreateBudgetExpense{
		Repository:     mockedRepository,
		EventPublisher: mockedPublisher,
	}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aBudgetExpense := BudgetExpense{
		Date:   *aDate,
		Amount: anAmount,
		Note:   "A_NOTE",
		Tags:   []tags.SearchTag{},
	}

	ctx := testutils.NewUserContext()

	mockedRepository.On("Save", ctx, &aBudgetExpense).Return(nil)
	mockedPublisher.On("CreateBudgetExpense", ctx, mock.Anything).Return(nil)

	err := uut.Execute(ctx, &aBudgetExpense)

	assert.Equal(t, nil, err)
	assert.Equal(t, []tags.SearchTag{{Key: "UNKNOWN", Value: "UNKNOWN"}}, aBudgetExpense.Tags)
}

func TestWhenANewBudgetExpenseCreationFails(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	uut := CreateBudgetExpense{
		Repository: mockedRepository,
		Logger:     testLogger,
	}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aBudgetExpense := BudgetExpense{
		Date:   *aDate,
		Amount: anAmount,
		Note:   "A_NOTE",
		Tags:   []tags.SearchTag{{Key: "super-market", Value: "super-market"}},
	}

	ctx := testutils.NewUserContext()

	saveError := errors.New("Budget Expense Save operation fails")
	mockedRepository.On("Save", ctx, &aBudgetExpense).Return(saveError)

	err := uut.Execute(ctx, &aBudgetExpense)

	assert.Equal(t, saveError, err)
	mockedRepository.AssertCalled(t, "Save", ctx, &aBudgetExpense)
}
