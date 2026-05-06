package attachment

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestSaveAttachmentSetsOwnerAndDelegatesToRepository(t *testing.T) {
	repo := new(AttachmentRepositoryMock)
	uut := SaveAttachment{repository: repo}

	ctx := testutils.NewUserContext()
	att := &Attachment{
		AttachmentMetadata: AttachmentMetadata{
			BudgetId:    "budget-123",
			FineName:    "receipt.pdf",
			ContentType: "application/pdf",
		},
		Content: []byte("content"),
	}

	repo.On("SaveAttachment", ctx, att).Return(nil)

	err := uut.Execute(ctx, att)

	assert.Equal(t, nil, err)
	assert.Equal(t, "A_USER_NAME", att.Owner)
	repo.AssertCalled(t, "SaveAttachment", ctx, att)
}

func TestSaveAttachmentPropagatesRepositoryError(t *testing.T) {
	repo := new(AttachmentRepositoryMock)
	uut := SaveAttachment{repository: repo}

	ctx := testutils.NewUserContext()
	att := &Attachment{
		AttachmentMetadata: AttachmentMetadata{BudgetId: "budget-123"},
	}
	repoErr := errors.New("storage unavailable")

	repo.On("SaveAttachment", ctx, att).Return(repoErr)

	err := uut.Execute(ctx, att)

	assert.Equal(t, repoErr, err)
	repo.AssertCalled(t, "SaveAttachment", ctx, att)
}

func TestSaveAttachmentReturnsErrorWhenNoUserInContext(t *testing.T) {
	repo := new(AttachmentRepositoryMock)
	uut := SaveAttachment{repository: repo}

	err := uut.Execute(context.Background(), &Attachment{})

	assert.NotEqual(t, nil, err)
	repo.AssertNotCalled(t, "SaveAttachment")
}
