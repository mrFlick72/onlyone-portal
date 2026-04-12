package revenue

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestWhenRevenuesAreFoundByYear(t *testing.T) {
	mockedRepository := new(RevenueRepositoryMock)
	uut := FindRevenue{Repository: mockedRepository}

	ctx := testutils.NewUserContext()
	year := date.NewYear(2018)

	start, _ := date.FirstDateOfMonth(date.JANUARY(), year)
	end, _ := date.LastDateOfMonth(date.DECEMBER(), year)

	aDate, _ := date.IsoDateFor("2018-06-15")
	amount, _ := money.MoneyFor("500.00")
	expected := []Revenue{
		{Id: "ID", UserName: "A_USER_NAME", Date: *aDate, Amount: amount, Note: "N"},
	}

	mockedRepository.On("FindByDateRange", ctx, start, end).Return(expected, nil)

	actual, err := uut.Execute(ctx, year)

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, actual)
	mockedRepository.AssertCalled(t, "FindByDateRange", ctx, start, end)
}
