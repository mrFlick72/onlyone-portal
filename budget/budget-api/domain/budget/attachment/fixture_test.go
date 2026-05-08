package attachment

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type AttachmentRepositoryMock struct {
	mock.Mock
}

func (m *AttachmentRepositoryMock) SaveAttachment(ctx context.Context, attachment *Attachment) error {
	args := m.Called(ctx, attachment)
	return args.Error(0)
}

func (m *AttachmentRepositoryMock) GenAttachment(ctx context.Context, attachmentId string) (*Attachment, error) {
	args := m.Called(ctx, attachmentId)
	if v := args.Get(0); v != nil {
		return v.(*Attachment), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AttachmentRepositoryMock) FindAllAttachment(ctx context.Context, budgetId string, budgetType BudgetType) ([]AttachmentMetadata, error) {
	args := m.Called(ctx, budgetId, budgetType)
	if v := args.Get(0); v != nil {
		return v.([]AttachmentMetadata), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AttachmentRepositoryMock) DeleteAttachment(ctx context.Context, attachmentId string) error {
	args := m.Called(ctx, attachmentId)
	return args.Error(0)
}
