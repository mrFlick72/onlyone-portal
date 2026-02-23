package expense

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestWhenABudgetExpenseDeletionSucceed(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	uut := DeleteBudgetExpense{
		Repository: mockedRepository,
	}

	ctx := testutils.NewUserContext()
	foundBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
	}
	mockedRepository.On("FindFor", &ctx, "A_BUDGET_ID").Return(foundBudgetExpense, nil)
	mockedRepository.On("Delete", &ctx, "A_BUDGET_ID").Return(nil)

	err := uut.Execute(&ctx, "A_BUDGET_ID")

	assert.Equal(t, nil, err)

	mockedRepository.AssertCalled(t, "FindFor", &ctx, "A_BUDGET_ID")
	mockedRepository.AssertCalled(t, "Delete", &ctx, "A_BUDGET_ID")
}

func TestWhenABudgetExpenseDeletionFails(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	uut := DeleteBudgetExpense{
		Repository: mockedRepository,
	}

	ctx := testutils.NewUserContext()

	foundBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_USER_NAME",
	}
	mockedRepository.On("FindFor", &ctx, "A_BUDGET_ID").Return(foundBudgetExpense, nil)
	mockedRepository.On("Delete", &ctx, "A_BUDGET_ID").Return(fmt.Errorf("budget expense delete with the id %s failed", "A_BUDGET_ID"))

	err := uut.Execute(&ctx, "A_BUDGET_ID")

	assert.NotEqual(t, nil, err)

	mockedRepository.AssertCalled(t, "FindFor", &ctx, "A_BUDGET_ID")
	mockedRepository.AssertCalled(t, "Delete", &ctx, "A_BUDGET_ID")
}

func TestWhenABudgetExpenseDeleteDoesDoneNothingBecauseTheBudgetExpenseDoesNotExist(t *testing.T) {

	mockedRepository := new(BudgetExpenseRepositoryMock)
	uut := DeleteBudgetExpense{
		Repository: mockedRepository,
	}
	ctx := testutils.NewUserContext()

	mockedRepository.On("FindFor", &ctx, "A_BUDGET_ID").Return(nil, fmt.Errorf("budget expense with the id %s was not found", "A_BUDGET_ID"))

	err := uut.Execute(&ctx, "A_BUDGET_ID")

	assert.NotEqual(t, nil, err)

	mockedRepository.AssertCalled(t, "FindFor", &ctx, "A_BUDGET_ID")
	mockedRepository.AssertNotCalled(t, "Delete", &ctx, "A_BUDGET_ID")
}

func TestWhenABudgetExpenseDeleteFailsBecauseUserOwnership(t *testing.T) {
mockedRepository := new(BudgetExpenseRepositoryMock)
	uut := DeleteBudgetExpense{
		Repository: mockedRepository,
	}

	ctx := testutils.NewUserContext()

	foundBudgetExpense := &BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "A_DIFFERENT_USER_NAME",
	}
	mockedRepository.On("FindFor", &ctx, "A_BUDGET_ID").Return(foundBudgetExpense, nil)

	err := uut.Execute(&ctx, "A_BUDGET_ID")

	assert.NotEqual(t, nil, err)

	mockedRepository.AssertCalled(t, "FindFor", &ctx, "A_BUDGET_ID")
	mockedRepository.AssertNotCalled(t, "Delete", &ctx, "A_BUDGET_ID")
}
