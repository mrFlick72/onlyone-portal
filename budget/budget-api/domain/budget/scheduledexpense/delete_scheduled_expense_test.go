package scheduledexpense

import (
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestWhenAnExistingScheduledExpenseIsDeleted(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := DeleteScheduledExpense{Repository: mockedRepository}

	existing := &ScheduledExpense{Id: "AN_ID", UserName: "A_USER_NAME", Status: StatusActive}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(existing, nil)
	mockedRepository.On("Delete", ctx, ScheduledExpenseId("AN_ID")).Return(nil)

	err := uut.Execute(ctx, "AN_ID")

	assert.Equal(t, nil, err)
	mockedRepository.AssertExpectations(t)
}

// An id that isn't in the caller's own partition — missing, or owned by
// someone else (indistinguishable by construction) — is rejected before the
// repository's Delete is ever reached.
func TestWhenDeletingAScheduledExpenseThatDoesNotExist(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := DeleteScheduledExpense{Repository: mockedRepository}

	ctx := testutils.NewUserContext()
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("MISSING_ID")).Return(nil, nil)

	err := uut.Execute(ctx, "MISSING_ID")

	assert.Equal(t, ErrScheduledExpenseNotFound, err)
	mockedRepository.AssertNotCalled(t, "Delete")
}

func TestWhenDeletingAScheduledExpenseFindForFails(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := DeleteScheduledExpense{Repository: mockedRepository}

	ctx := testutils.NewUserContext()
	findError := errors.New("ScheduledExpense FindFor operation fails")
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(nil, findError)

	err := uut.Execute(ctx, "AN_ID")

	assert.Equal(t, findError, err)
	mockedRepository.AssertNotCalled(t, "Delete")
}

func TestWhenDeletingAScheduledExpenseDeleteFails(t *testing.T) {
	mockedRepository := new(ScheduledExpenseRepositoryMock)
	uut := DeleteScheduledExpense{Repository: mockedRepository}

	existing := &ScheduledExpense{Id: "AN_ID", UserName: "A_USER_NAME", Status: StatusActive}

	ctx := testutils.NewUserContext()
	deleteError := errors.New("ScheduledExpense Delete operation fails")
	mockedRepository.On("FindFor", ctx, ScheduledExpenseId("AN_ID")).Return(existing, nil)
	mockedRepository.On("Delete", ctx, ScheduledExpenseId("AN_ID")).Return(deleteError)

	err := uut.Execute(ctx, "AN_ID")

	assert.Equal(t, deleteError, err)
}
