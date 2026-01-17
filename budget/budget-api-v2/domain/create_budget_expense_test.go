package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/stretchr/testify/mock"
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

	UserName := "A_USER_NAME"
	user := security.User{UserName: &UserName, Authorities: nil}
	ctx := context.WithValue(context.TODO(), "user", user)

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

	UserName := "A_USER_NAME"
	user := security.User{UserName: &UserName, Authorities: nil}
	ctx := context.WithValue(context.TODO(), "user", user)

	saveError := errors.New("Budget Expense Save operation fails")
	mockedRepository.On("Save", &ctx, &aBudgetExpense).Return(saveError)

	err := uut.Execute(&ctx, &aBudgetExpense)

	assert.Equal(t, saveError, err)
	mockedRepository.AssertCalled(t, "Save", &ctx, &aBudgetExpense)
}

type BudgetExpenseRepositoryMock struct {
	mock.Mock
}

func (mock *BudgetExpenseRepositoryMock) FindFor(ctx *context.Context, budgetExpenseId BudgetExpenseId) (*BudgetExpense, error) {
	args := mock.Called(ctx, budgetExpenseId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BudgetExpense), args.Error(1)
}

func (mock *BudgetExpenseRepositoryMock) FindByDateRange(ctx *context.Context, userName UserName, star Date, end Date, searchTags []string) (*[]BudgetExpense, error) {
	return nil, nil
}
func (mock *BudgetExpenseRepositoryMock) Save(ctx *context.Context, budgetExpense *BudgetExpense) error {
	budgetExpense.Id = "A_BUDGET_ID"
	args := mock.Called(ctx, budgetExpense)
	return args.Error(0)

}
func (mock *BudgetExpenseRepositoryMock) Delete(ctx *context.Context, idBudgetExpense BudgetExpenseId) error {
	args := mock.Called(ctx, idBudgetExpense)
	return args.Error(0)
}
