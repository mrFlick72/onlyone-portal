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

func (m *AttachmentRepositoryMock) GenAttachment(ctx context.Context, attachmentId string) (error, *Attachment) {
	args := m.Called(ctx, attachmentId)
	if v := args.Get(1); v != nil {
		return args.Error(0), v.(*Attachment)
	}
	return args.Error(0), nil
}

func (m *AttachmentRepositoryMock) FindAllAttachment(ctx context.Context, budgetId string) (error, []AttachmentMetadata) {
	args := m.Called(ctx, budgetId)
	if v := args.Get(1); v != nil {
		return args.Error(0), v.([]AttachmentMetadata)
	}
	return args.Error(0), nil
}

func (m *AttachmentRepositoryMock) DeleteAttachment(ctx context.Context, attachmentId string) error {
	args := m.Called(ctx, attachmentId)
	return args.Error(0)
}
