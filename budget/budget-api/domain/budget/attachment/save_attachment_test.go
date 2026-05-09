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

func TestGetAttachmentGetOneByDelegatesToRepository(t *testing.T) {
	repo := new(AttachmentRepositoryMock)
	uut := NewGetAttachment(repo)

	ctx := testutils.NewUserContext()
	expected := &Attachment{
		AttachmentMetadata: AttachmentMetadata{AttachmentId: "att-1", FineName: "f.pdf"},
		Content:            []byte("c"),
	}
	repo.On("GenAttachment", ctx, "att-1").Return(expected, nil)

	got, err := uut.GetOneBy(ctx, "att-1")

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, got)
	repo.AssertCalled(t, "GenAttachment", ctx, "att-1")
}

func TestGetAttachmentGetOneByPropagatesError(t *testing.T) {
	repo := new(AttachmentRepositoryMock)
	uut := NewGetAttachment(repo)

	ctx := testutils.NewUserContext()
	repoErr := errors.New("not found")
	repo.On("GenAttachment", ctx, "missing").Return(nil, repoErr)

	got, err := uut.GetOneBy(ctx, "missing")

	assert.Equal(t, repoErr, err)
	assert.Equal(t, (*Attachment)(nil), got)
}

func TestGetAttachmentGetAllByDelegatesToRepository(t *testing.T) {
	repo := new(AttachmentRepositoryMock)
	uut := NewGetAttachment(repo)

	ctx := testutils.NewUserContext()
	expected := []AttachmentMetadata{
		{AttachmentId: "a", BudgetId: "b", BudgetType: "expense"},
	}
	repo.On("FindAllAttachment", ctx, "b", Expense).Return(expected, nil)

	got, err := uut.GetAllBy(ctx, "b", Expense)

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, got)
	repo.AssertCalled(t, "FindAllAttachment", ctx, "b", Expense)
}

func TestActionFacadeRoutesEachMethodToCorrespondingAction(t *testing.T) {
	repo := new(AttachmentRepositoryMock)
	facade := AttachmentActionsFacade{
		SaveAttachmentAction: NewSaveAttachment(repo),
		GetAttachmentAction:  NewGetAttachment(repo),
	}

	ctx := testutils.NewUserContext()
	att := &Attachment{AttachmentMetadata: AttachmentMetadata{BudgetId: "b"}}
	repo.On("SaveAttachment", ctx, att).Return(nil)
	assert.Equal(t, nil, facade.SaveAttachment(ctx, att))

	expected := &Attachment{AttachmentMetadata: AttachmentMetadata{AttachmentId: "x"}}
	repo.On("GenAttachment", ctx, "x").Return(expected, nil)
	got, err := facade.GetAttachmentBy(ctx, "x")
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, got)

	list := []AttachmentMetadata{{AttachmentId: "a"}}
	repo.On("FindAllAttachment", ctx, "b", Revenue).Return(list, nil)
	gotList, err := facade.GetAttachments(ctx, "b", Revenue)
	assert.Equal(t, nil, err)
	assert.Equal(t, list, gotList)
}
