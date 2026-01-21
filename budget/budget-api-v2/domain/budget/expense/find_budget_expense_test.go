package expense

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
)

func TestBudgetExpenseTotalBySearchTagsWithConstraints(t *testing.T) {

	ctx := newUserContext()
	mockedSearchTagRepository := new(SearchTagRepositoryMock)
	mockedBudgetExpenseRepository := new(BudgetExpenseRepositoryMock)
	uut := FindSpentBudget{
		budgetExpenseRepository: mockedBudgetExpenseRepository,
		searchTagRepository:     mockedSearchTagRepository,
	}

	mockedSearchTagRepository.On("GetTagBy", &ctx, "super-market").Return(&tags.SearchTag{Key: "super-market", Value: "super-market"}, nil)
        mockedSearchTagRepository.On("GetTagBy", &ctx, "dinner").Return(&tags.SearchTag{Key: "dinner", Value: "dinner"}, nil)
        
	mockedBudgetExpenseRepository.On("FindByDateRange", &ctx, "A_USER_NAME", safeDateFor("01/04/2018"), safeDateFor("30/04/2018"), []string{"dinner", "super-market"}).
		Return(&[]BudgetExpense{
			{Id: "1", UserName: "A_USER_NAME", Date: safeDateFor("15/04/2018"), Amount: safeMoneyFor("10.00"), Note: "dinner", Tag: "dinner"},
			{Id: "2", UserName: "A_USER_NAME", Date: safeDateFor("01/04/2018"), Amount: safeMoneyFor("12.50"), Note: "super-market", Tag: "super-market"},
			{Id: "3", UserName: "A_USER_NAME", Date: safeDateFor("05/04/2018"), Amount: safeMoneyFor("12.50"), Note: "super-market", Tag: "super-market"},
			{Id: "4", UserName: "A_USER_NAME", Date: safeDateFor("04/04/2018"), Amount: safeMoneyFor("20.00"), Note: "dinner", Tag: "dinner"},
			{Id: "5", UserName: "A_USER_NAME", Date: safeDateFor("03/04/2018"), Amount: safeMoneyFor("12.50"), Note: "super-market", Tag: "super-market"},
			{Id: "6", UserName: "A_USER_NAME", Date: safeDateFor("02/04/2018"), Amount: safeMoneyFor("12.50"), Note: "super-market", Tag: "super-market"},
			{Id: "7", UserName: "A_USER_NAME", Date: safeDateFor("01/04/2018"), Amount: safeMoneyFor("15.00"), Note: "dinner", Tag: "dinner"},
		}, nil)
	actual, err := uut.Execute(&ctx, date.APRIL(), date.NewYear(2018), []string{"dinner", "super-market"})

	assert.NotEqual(t, nil, actual)
	assert.Equal(t, nil, err)

    mockedSearchTagRepository.AssertCalled(t, "GetTagBy", &ctx, "dinner")
	mockedSearchTagRepository.AssertCalled(t, "GetTagBy", &ctx, "super-market")
	mockedBudgetExpenseRepository.AssertCalled(t, "FindByDateRange", &ctx, "A_USER_NAME", safeDateFor("01/04/2018"), safeDateFor("30/04/2018"), []string{"dinner", "super-market"})
}

func TestBudgetExpenseTotalBySearchTagsWithoutConstraints(t *testing.T) {
		ctx := newUserContext()
	mockedSearchTagRepository := new(SearchTagRepositoryMock)
	mockedBudgetExpenseRepository := new(BudgetExpenseRepositoryMock)
	uut := FindSpentBudget{
		budgetExpenseRepository: mockedBudgetExpenseRepository,
		searchTagRepository:     mockedSearchTagRepository,
	}

	mockedSearchTagRepository.On("GetTagBy", &ctx, "super-market").Return(&tags.SearchTag{Key: "super-market", Value: "super-market"}, nil)
        mockedSearchTagRepository.On("GetTagBy", &ctx, "dinner").Return(&tags.SearchTag{Key: "dinner", Value: "dinner"}, nil)

	mockedBudgetExpenseRepository.On("FindByDateRange", &ctx, "A_USER_NAME", safeDateFor("01/04/2018"), safeDateFor("30/04/2018"), []string{}).
		Return(&[]BudgetExpense{
			{Id: "1", UserName: "A_USER_NAME", Date: safeDateFor("15/04/2018"), Amount: safeMoneyFor("10.00"), Note: "dinner", Tag: "dinner"},
			{Id: "2", UserName: "A_USER_NAME", Date: safeDateFor("01/04/2018"), Amount: safeMoneyFor("12.50"), Note: "super-market", Tag: "super-market"},
			{Id: "3", UserName: "A_USER_NAME", Date: safeDateFor("05/04/2018"), Amount: safeMoneyFor("12.50"), Note: "super-market", Tag: "super-market"},
			{Id: "4", UserName: "A_USER_NAME", Date: safeDateFor("04/04/2018"), Amount: safeMoneyFor("20.00"), Note: "dinner", Tag: "dinner"},
			{Id: "5", UserName: "A_USER_NAME", Date: safeDateFor("03/04/2018"), Amount: safeMoneyFor("12.50"), Note: "super-market", Tag: "super-market"},
			{Id: "6", UserName: "A_USER_NAME", Date: safeDateFor("02/04/2018"), Amount: safeMoneyFor("12.50"), Note: "super-market", Tag: "super-market"},
			{Id: "7", UserName: "A_USER_NAME", Date: safeDateFor("01/04/2018"), Amount: safeMoneyFor("15.00"), Note: "dinner", Tag: "dinner"},
		}, nil)
	actual, err := uut.Execute(&ctx, date.APRIL(), date.NewYear(2018), []string{})

	assert.NotEqual(t, nil, actual)
	assert.Equal(t, nil, err)

    mockedSearchTagRepository.AssertCalled(t, "GetTagBy", &ctx, "dinner")
	mockedSearchTagRepository.AssertCalled(t, "GetTagBy", &ctx, "super-market")
	mockedBudgetExpenseRepository.AssertCalled(t, "FindByDateRange", &ctx, "A_USER_NAME", safeDateFor("01/04/2018"), safeDateFor("30/04/2018"), []string{})
}
