package attachment

import (
	"context"
)

type AttachmentActions interface {
	SaveAttachment(ctx context.Context, attachment *Attachment) error
	GetAttachmentBy(ctx context.Context, attachmentId string) (*Attachment, error)
	GetAttachments(ctx context.Context, budgetId string) ([]AttachmentMetadata, error)
	DeleteAttachment(ctx context.Context, attachmentId string) error
}

type AttachmentActionsFacade struct {
}

func (facade *AttachmentActionsFacade) AddAttachment(ctx context.Context, attachment *Attachment) error {
	panic("TODO")
}

func (facade *AttachmentActionsFacade) GetAttachmentBy(ctx context.Context, attachmentId string) error {
	panic("TODO")
}

func (facade *AttachmentActionsFacade) GetAttachments(ctx context.Context, budgetId string) (*Attachment, error) {
	panic("TODO")
}

func (facade *AttachmentActionsFacade) DeleteAttachment(ctx context.Context, attachmentId string) error {
	panic("TODO")
}
