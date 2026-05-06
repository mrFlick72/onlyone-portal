package attachment

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/stretchr/testify/mock"
)

type AttachmentActionsMock struct {
	mock.Mock
}

func (m *AttachmentActionsMock) SaveAttachment(ctx context.Context, att *attachment.Attachment) error {
	args := m.Called(ctx, att)
	return args.Error(0)
}

func (m *AttachmentActionsMock) GetAttachmentBy(ctx context.Context, attachmentId string) (*attachment.Attachment, error) {
	args := m.Called(ctx, attachmentId)
	if v := args.Get(0); v != nil {
		return v.(*attachment.Attachment), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AttachmentActionsMock) GetAttachments(ctx context.Context, budgetId string) ([]attachment.AttachmentMetadata, error) {
	args := m.Called(ctx, budgetId)
	if v := args.Get(0); v != nil {
		return v.([]attachment.AttachmentMetadata), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AttachmentActionsMock) DeleteAttachment(ctx context.Context, attachmentId string) error {
	args := m.Called(ctx, attachmentId)
	return args.Error(0)
}

type ContextFactoryConverterMock struct {
	mock.Mock
}

func (m *ContextFactoryConverterMock) CreateContextFromGin(c *gin.Context) context.Context {
	args := m.Called(c)
	return args.Get(0).(context.Context)
}
