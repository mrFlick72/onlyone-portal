package expense

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestBudgetExpenseTotalBySearchTagsWithConstraints(t *testing.T) {

	ctx := testutils.NewUserContext()
	mockedSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockedBudgetExpenseRepository := new(BudgetExpenseRepositoryMock)
	uut := FindSpentBudget{
		BudgetExpenseRepository: mockedBudgetExpenseRepository,
		SearchTagRepository:     mockedSearchTagRepository,
	}

	mockedSearchTagRepository.On("GetTagBy", ctx, "super-market").Return(&tags.SearchTag{Key: "super-market", Value: "super-market"}, nil)
	mockedSearchTagRepository.On("GetTagBy", ctx, "dinner").Return(&tags.SearchTag{Key: "dinner", Value: "dinner"}, nil)

	mockedBudgetExpenseRepository.On("FindByDateRange", ctx, testutils.SafeDateFor("01/04/2018"), testutils.SafeDateFor("30/04/2018"), []string{"dinner", "super-market"}).
		Return([]BudgetExpense{
			{Id: "1", UserName: "A_USER_NAME", Date: testutils.SafeDateFor("15/04/2018"), Amount: testutils.SafeMoneyFor("10.00"), Note: "dinner", Tags: []tags.SearchTag{{Key: "dinner", Value: "dinner"}}},
			{Id: "2", UserName: "A_USER_NAME", Date: testutils.SafeDateFor("01/04/2018"), Amount: testutils.SafeMoneyFor("12.50"), Note: "super-market", Tags: []tags.SearchTag{{Key: "super-market", Value: "super-market"}}},
			{Id: "3", UserName: "A_USER_NAME", Date: testutils.SafeDateFor("05/04/2018"), Amount: testutils.SafeMoneyFor("12.50"), Note: "super-market", Tags: []tags.SearchTag{{Key: "super-market", Value: "super-market"}}},
			{Id: "4", UserName: "A_USER_NAME", Date: testutils.SafeDateFor("04/04/2018"), Amount: testutils.SafeMoneyFor("20.00"), Note: "dinner", Tags: []tags.SearchTag{{Key: "dinner", Value: "dinner"}}},
			{Id: "5", UserName: "A_USER_NAME", Date: testutils.SafeDateFor("03/04/2018"), Amount: testutils.SafeMoneyFor("12.50"), Note: "super-market", Tags: []tags.SearchTag{{Key: "super-market", Value: "super-market"}}},
			{Id: "6", UserName: "A_USER_NAME", Date: testutils.SafeDateFor("02/04/2018"), Amount: testutils.SafeMoneyFor("12.50"), Note: "super-market", Tags: []tags.SearchTag{{Key: "super-market", Value: "super-market"}}},
			{Id: "7", UserName: "A_USER_NAME", Date: testutils.SafeDateFor("01/04/2018"), Amount: testutils.SafeMoneyFor("15.00"), Note: "dinner", Tags: []tags.SearchTag{{Key: "dinner", Value: "dinner"}}},
		}, nil)
	actual, err := uut.Execute(ctx, date.APRIL(), date.NewYear(2018), []string{"dinner", "super-market"})

	assert.NotEqual(t, nil, actual)
	assert.Equal(t, nil, err)

	mockedSearchTagRepository.AssertCalled(t, "GetTagBy", ctx, "dinner")
	mockedSearchTagRepository.AssertCalled(t, "GetTagBy", ctx, "super-market")
	mockedBudgetExpenseRepository.AssertCalled(t, "FindByDateRange", ctx, testutils.SafeDateFor("01/04/2018"), testutils.SafeDateFor("30/04/2018"), []string{"dinner", "super-market"})
}

// Regression for review finding #2: an expense that carries the same resolved
// tag key more than once (e.g. two deleted tags both resolving to UNKNOWN) must
// contribute its amount to that tag's total only once; an expense tagged with
// distinct keys still contributes to each. Before the fix, [UNKNOWN, UNKNOWN]
// added 10.00 twice, making the UNKNOWN total 25.00 instead of 15.00.
func TestTotalForSearchTagsCountsDuplicateResolvedTagsOncePerExpense(t *testing.T) {
	unknown := tags.UnknownSentinel()
	food := tags.SearchTag{Key: "food", Value: "Food"}
	transport := tags.SearchTag{Key: "transport", Value: "Transport"}

	spentBudget := NewSpentBudget(
		[]BudgetExpense{
			{Id: "1", Amount: testutils.SafeMoneyFor("10.00"), Tags: []tags.SearchTag{unknown, unknown}},
			{Id: "2", Amount: testutils.SafeMoneyFor("5.00"), Tags: []tags.SearchTag{unknown}},
			{Id: "3", Amount: testutils.SafeMoneyFor("7.00"), Tags: []tags.SearchTag{food, transport}},
		},
		[]tags.SearchTag{unknown, food, transport},
	)

	totals := spentBudget.TotalForSearchTags()

	assert.Equal(t, "15.00", totals[unknown].StringifyAmount())
	assert.Equal(t, "7.00", totals[food].StringifyAmount())
	assert.Equal(t, "7.00", totals[transport].StringifyAmount())
}

func TestBudgetExpenseTotalBySearchTagsWithoutConstraints(t *testing.T) {
	ctx := testutils.NewUserContext()
	mockedSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockedBudgetExpenseRepository := new(BudgetExpenseRepositoryMock)
	uut := FindSpentBudget{
		BudgetExpenseRepository: mockedBudgetExpenseRepository,
		SearchTagRepository:     mockedSearchTagRepository,
	}

	mockedSearchTagRepository.On("GetTagBy", ctx, "super-market").Return(&tags.SearchTag{Key: "super-market", Value: "super-market"}, nil)
	mockedSearchTagRepository.On("GetTagBy", ctx, "dinner").Return(&tags.SearchTag{Key: "dinner", Value: "dinner"}, nil)

	mockedBudgetExpenseRepository.On("FindByDateRange", ctx, testutils.SafeDateFor("01/04/2018"), testutils.SafeDateFor("30/04/2018"), []string{}).
		Return([]BudgetExpense{
			{Id: "1", UserName: "A_USER_NAME", Date: testutils.SafeDateFor("15/04/2018"), Amount: testutils.SafeMoneyFor("10.00"), Note: "dinner", Tags: []tags.SearchTag{{Key: "dinner", Value: "dinner"}}},
			{Id: "2", UserName: "A_USER_NAME", Date: testutils.SafeDateFor("01/04/2018"), Amount: testutils.SafeMoneyFor("12.50"), Note: "super-market", Tags: []tags.SearchTag{{Key: "super-market", Value: "super-market"}}},
			{Id: "3", UserName: "A_USER_NAME", Date: testutils.SafeDateFor("05/04/2018"), Amount: testutils.SafeMoneyFor("12.50"), Note: "super-market", Tags: []tags.SearchTag{{Key: "super-market", Value: "super-market"}}},
			{Id: "4", UserName: "A_USER_NAME", Date: testutils.SafeDateFor("04/04/2018"), Amount: testutils.SafeMoneyFor("20.00"), Note: "dinner", Tags: []tags.SearchTag{{Key: "dinner", Value: "dinner"}}},
			{Id: "5", UserName: "A_USER_NAME", Date: testutils.SafeDateFor("03/04/2018"), Amount: testutils.SafeMoneyFor("12.50"), Note: "super-market", Tags: []tags.SearchTag{{Key: "super-market", Value: "super-market"}}},
			{Id: "6", UserName: "A_USER_NAME", Date: testutils.SafeDateFor("02/04/2018"), Amount: testutils.SafeMoneyFor("12.50"), Note: "super-market", Tags: []tags.SearchTag{{Key: "super-market", Value: "super-market"}}},
			{Id: "7", UserName: "A_USER_NAME", Date: testutils.SafeDateFor("01/04/2018"), Amount: testutils.SafeMoneyFor("15.00"), Note: "dinner", Tags: []tags.SearchTag{{Key: "dinner", Value: "dinner"}}},
		}, nil)
	actual, err := uut.Execute(ctx, date.APRIL(), date.NewYear(2018), []string{})

	assert.NotEqual(t, nil, actual)
	assert.Equal(t, nil, err)

	mockedSearchTagRepository.AssertCalled(t, "GetTagBy", ctx, "dinner")
	mockedSearchTagRepository.AssertCalled(t, "GetTagBy", ctx, "super-market")
	mockedBudgetExpenseRepository.AssertCalled(t, "FindByDateRange", ctx, testutils.SafeDateFor("01/04/2018"), testutils.SafeDateFor("30/04/2018"), []string{})
}
