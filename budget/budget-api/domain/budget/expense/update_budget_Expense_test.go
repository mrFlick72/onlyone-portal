package expense

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/stretchr/testify/mock"
)

func TestWhenABudgetExpenseUpdateSucceed(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	mockedPublisher := new(BudgetExpenseEventPublisherMock)
	uut := UpdateBudgetExpense{
		Repository:     mockedRepository,
		EventPublisher: mockedPublisher,
	}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aBudgetExpense := BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
		Date:     *aDate,
		Amount:   anAmount,
		Note:     "A_NOTE",
		Tags:     []tags.SearchTag{{Key: "super-market", Value: "super-market"}},
	}

	ctx := testutils.NewUserContext()

	foundBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
	}
	mockedRepository.On("FindFor", ctx, "A_BUDGET_ID").Return(foundBudgetExpense, nil)
	mockedRepository.On("Save", ctx, &aBudgetExpense).Return(nil)
	mockedPublisher.On("UpdateBudgetExpense", ctx, mock.Anything).Return(nil)

	err := uut.Execute(ctx, &aBudgetExpense)

	assert.Equal(t, nil, err)

	mockedRepository.AssertCalled(t, "FindFor", ctx, "A_BUDGET_ID")
	mockedRepository.AssertCalled(t, "Save", ctx, &aBudgetExpense)
	mockedPublisher.AssertCalled(t, "UpdateBudgetExpense", ctx, aBudgetExpense)
}
func TestWhenABudgetExpenseUpdateWithNoTagsDefaultsToUnknown(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	mockedPublisher := new(BudgetExpenseEventPublisherMock)
	uut := UpdateBudgetExpense{
		Repository:     mockedRepository,
		EventPublisher: mockedPublisher,
	}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aBudgetExpense := BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
		Date:     *aDate,
		Amount:   anAmount,
		Note:     "A_NOTE",
		Tags:     []tags.SearchTag{},
	}

	ctx := testutils.NewUserContext()

	foundBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
	}
	mockedRepository.On("FindFor", ctx, "A_BUDGET_ID").Return(foundBudgetExpense, nil)
	mockedRepository.On("Save", ctx, &aBudgetExpense).Return(nil)
	mockedPublisher.On("UpdateBudgetExpense", ctx, mock.Anything).Return(nil)

	err := uut.Execute(ctx, &aBudgetExpense)

	assert.Equal(t, nil, err)
	assert.Equal(t, []tags.SearchTag{{Key: "UNKNOWN", Value: "UNKNOWN"}}, aBudgetExpense.Tags)
}

func TestWhenABudgetExpenseUpdateDoesDoneNothingBecauseTheBudgetExpenseDoesNotExist(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	uut := UpdateBudgetExpense{
		Repository: mockedRepository,
	}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aBudgetExpense := BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
		Date:     *aDate,
		Amount:   anAmount,
		Note:     "A_NOTE",
		Tags:     []tags.SearchTag{{Key: "super-market", Value: "super-market"}},
	}

	ctx := testutils.NewUserContext()

	mockedRepository.On("FindFor", ctx, "A_BUDGET_ID").Return(nil, fmt.Errorf("budget expense with the id %s was not found", "A_BUDGET_ID"))

	err := uut.Execute(ctx, &aBudgetExpense)

	assert.NotEqual(t, nil, err)

	mockedRepository.AssertCalled(t, "FindFor", ctx, "A_BUDGET_ID")
	mockedRepository.AssertNotCalled(t, "Save", ctx, &aBudgetExpense)

}
func TestWhenABudgetExpenseUpdateFails(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	uut := UpdateBudgetExpense{
		Repository: mockedRepository,
	}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aBudgetExpense := BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
		Date:     *aDate,
		Amount:   anAmount,
		Note:     "A_NOTE",
		Tags:     []tags.SearchTag{{Key: "super-market", Value: "super-market"}},
	}

	ctx := testutils.NewUserContext()

	foundBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
	}
	mockedRepository.On("FindFor", ctx, "A_BUDGET_ID").Return(foundBudgetExpense, nil)
	mockedRepository.On("Save", ctx, &aBudgetExpense).Return(fmt.Errorf("budget expense save with the id %s failed", "A_BUDGET_ID"))

	err := uut.Execute(ctx, &aBudgetExpense)

	assert.NotEqual(t, nil, err)

	mockedRepository.AssertCalled(t, "FindFor", ctx, "A_BUDGET_ID")
	mockedRepository.AssertCalled(t, "Save", ctx, &aBudgetExpense)
}

func TestWhenABudgetExpenseUpdateFailsBecauseUserOwnership(t *testing.T) {
	mockedRepository := new(BudgetExpenseRepositoryMock)
	uut := UpdateBudgetExpense{
		Repository: mockedRepository,
	}

	ctx := testutils.NewUserContext()
	foundBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_DIFFERENT_USER_NAME",
	}
	aBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
	}
	mockedRepository.On("FindFor", ctx, "A_BUDGET_ID").Return(foundBudgetExpense, nil)

	err := uut.Execute(ctx, aBudgetExpense)

	assert.NotEqual(t, nil, err)

	mockedRepository.AssertCalled(t, "FindFor", ctx, "A_BUDGET_ID")
	mockedRepository.AssertNotCalled(t, "Save", ctx, aBudgetExpense)
}
