package revenue

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestWhenARevenueDeletionSucceed(t *testing.T) {
	mockedRepository := new(RevenueRepositoryMock)
	uut := DeleteRevenue{Repository: mockedRepository}

	ctx := testutils.NewUserContext()
	found := &Revenue{Id: "A_REVENUE_ID", UserName: "A_USER_NAME"}
	mockedRepository.On("FindFor", ctx, "A_REVENUE_ID").Return(found, nil)
	mockedRepository.On("Delete", ctx, "A_REVENUE_ID").Return(nil)

	err := uut.Execute(ctx, "A_REVENUE_ID")

	assert.Equal(t, nil, err)
	mockedRepository.AssertCalled(t, "FindFor", ctx, "A_REVENUE_ID")
	mockedRepository.AssertCalled(t, "Delete", ctx, "A_REVENUE_ID")
}

func TestWhenARevenueDeletionFails(t *testing.T) {
	mockedRepository := new(RevenueRepositoryMock)
	uut := DeleteRevenue{Repository: mockedRepository}

	ctx := testutils.NewUserContext()
	found := &Revenue{Id: "A_REVENUE_ID", UserName: "A_USER_NAME"}
	mockedRepository.On("FindFor", ctx, "A_REVENUE_ID").Return(found, nil)
	mockedRepository.On("Delete", ctx, "A_REVENUE_ID").Return(fmt.Errorf("delete failed"))

	err := uut.Execute(ctx, "A_REVENUE_ID")

	assert.NotEqual(t, nil, err)
	mockedRepository.AssertCalled(t, "Delete", ctx, "A_REVENUE_ID")
}

func TestWhenARevenueDeleteDoesNothingBecauseNotFound(t *testing.T) {
	mockedRepository := new(RevenueRepositoryMock)
	uut := DeleteRevenue{Repository: mockedRepository}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, "A_REVENUE_ID").Return(nil, fmt.Errorf("not found"))

	err := uut.Execute(ctx, "A_REVENUE_ID")

	assert.NotEqual(t, nil, err)
	mockedRepository.AssertCalled(t, "FindFor", ctx, "A_REVENUE_ID")
	mockedRepository.AssertNotCalled(t, "Delete", ctx, "A_REVENUE_ID")
}

func TestWhenARevenueDeleteFailsBecauseUserOwnership(t *testing.T) {
	mockedRepository := new(RevenueRepositoryMock)
	uut := DeleteRevenue{Repository: mockedRepository}

	ctx := testutils.NewUserContext()
	found := &Revenue{Id: "A_REVENUE_ID", UserName: "ANOTHER_USER_NAME"}
	mockedRepository.On("FindFor", ctx, "A_REVENUE_ID").Return(found, nil)

	err := uut.Execute(ctx, "A_REVENUE_ID")

	assert.NotEqual(t, nil, err)
	mockedRepository.AssertCalled(t, "FindFor", ctx, "A_REVENUE_ID")
	mockedRepository.AssertNotCalled(t, "Delete", ctx, "A_REVENUE_ID")
}
