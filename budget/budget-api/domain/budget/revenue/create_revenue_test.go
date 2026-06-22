package revenue

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestWhenANewRevenueIsCreated(t *testing.T) {
	mockedRepository := new(RevenueRepositoryMock)
	uut := CreateRevenue{Repository: mockedRepository}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aRevenue := Revenue{
		Date:   *aDate,
		Amount: anAmount,
		Note:   "A_NOTE",
	}

	ctx := testutils.NewUserContext()
	mockedRepository.On("Save", ctx, &aRevenue).Return(nil)

	err := uut.Execute(ctx, &aRevenue)

	assert.Equal(t, nil, err)
	assert.Equal(t, "A_USER_NAME", aRevenue.UserName)
	assert.Equal(t, "A_REVENUE_ID", aRevenue.Id)
	mockedRepository.AssertCalled(t, "Save", ctx, &aRevenue)
}

func TestWhenANewRevenueHasNoTagsItDefaultsToUnknown(t *testing.T) {
	mockedRepository := new(RevenueRepositoryMock)
	uut := CreateRevenue{Repository: mockedRepository}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aRevenue := Revenue{Date: *aDate, Amount: anAmount, Note: "A_NOTE"}

	ctx := testutils.NewUserContext()
	mockedRepository.On("Save", ctx, &aRevenue).Return(nil)

	err := uut.Execute(ctx, &aRevenue)

	assert.Equal(t, nil, err)
	assert.Equal(t, []tags.SearchTag{tags.UnknownSentinel()}, aRevenue.Tags)
}

func TestWhenANewRevenueHasTagsTheyArePreserved(t *testing.T) {
	mockedRepository := new(RevenueRepositoryMock)
	uut := CreateRevenue{Repository: mockedRepository}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	provided := []tags.SearchTag{{Key: "salary", Value: "Salary"}}
	aRevenue := Revenue{Date: *aDate, Amount: anAmount, Note: "A_NOTE", Tags: provided}

	ctx := testutils.NewUserContext()
	mockedRepository.On("Save", ctx, &aRevenue).Return(nil)

	err := uut.Execute(ctx, &aRevenue)

	assert.Equal(t, nil, err)
	assert.Equal(t, provided, aRevenue.Tags)
}

func TestWhenANewRevenueCreationFailsBecauseNoUserInContext(t *testing.T) {
	mockedRepository := new(RevenueRepositoryMock)
	uut := CreateRevenue{Repository: mockedRepository}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aRevenue := Revenue{
		Date:   *aDate,
		Amount: anAmount,
		Note:   "A_NOTE",
	}

	err := uut.Execute(context.Background(), &aRevenue)

	assert.NotEqual(t, nil, err)
	mockedRepository.AssertNotCalled(t, "Save")
}

func TestWhenANewRevenueCreationFails(t *testing.T) {
	mockedRepository := new(RevenueRepositoryMock)
	uut := CreateRevenue{Repository: mockedRepository}

	aDate, _ := date.IsoDateFor("2018-01-01")
	anAmount, _ := money.MoneyFor("1.00")
	aRevenue := Revenue{
		Date:   *aDate,
		Amount: anAmount,
		Note:   "A_NOTE",
	}

	ctx := testutils.NewUserContext()
	saveError := errors.New("Revenue Save operation fails")
	mockedRepository.On("Save", ctx, &aRevenue).Return(saveError)

	err := uut.Execute(ctx, &aRevenue)

	assert.Equal(t, saveError, err)
	mockedRepository.AssertCalled(t, "Save", ctx, &aRevenue)
}
