package revenue

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestWhenARevenueIsUpdatedByTheOwner(t *testing.T) {
	mockedRepository := new(RevenueRepositoryMock)
	uut := UpdateRevenue{Repository: mockedRepository}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aRevenue := &Revenue{
		Id:     "A_REVENUE_ID",
		Date:   *aDate,
		Amount: anAmount,
		Note:   "A_NOTE",
	}

	ctx := testutils.NewUserContext()

	existing := &Revenue{Id: "A_REVENUE_ID", UserName: "A_USER_NAME"}
	mockedRepository.On("FindFor", ctx, "A_REVENUE_ID").Return(existing, nil)
	mockedRepository.On("Save", ctx, aRevenue).Return(nil)

	err := uut.Execute(ctx, aRevenue)

	assert.Equal(t, nil, err)
	mockedRepository.AssertCalled(t, "FindFor", ctx, "A_REVENUE_ID")
	mockedRepository.AssertCalled(t, "Save", ctx, aRevenue)
}

func TestWhenARevenueUpdateFailsBecauseOwnership(t *testing.T) {
	mockedRepository := new(RevenueRepositoryMock)
	uut := UpdateRevenue{Repository: mockedRepository}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aRevenue := &Revenue{
		Id:     "A_REVENUE_ID",
		Date:   *aDate,
		Amount: anAmount,
	}

	ctx := testutils.NewUserContext()

	existing := &Revenue{Id: "A_REVENUE_ID", UserName: "ANOTHER_USER"}
	mockedRepository.On("FindFor", ctx, "A_REVENUE_ID").Return(existing, nil)

	err := uut.Execute(ctx, aRevenue)

	assert.NotEqual(t, nil, err)
	mockedRepository.AssertCalled(t, "FindFor", ctx, "A_REVENUE_ID")
	mockedRepository.AssertNotCalled(t, "Save", ctx, aRevenue)
}
