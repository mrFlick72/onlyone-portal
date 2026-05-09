package attachment

import (
	"context"
)

type AttachmentActions interface {
	SaveAttachment(ctx context.Context, attachment *Attachment) error
	GetAttachmentBy(ctx context.Context, attachmentId string) (*Attachment, error)
	GetAttachments(ctx context.Context, budgetId string, budgetType BudgetType) ([]AttachmentMetadata, error)
	DeleteAttachment(ctx context.Context, attachmentId string) error
}

type AttachmentActionsFacade struct {
	SaveAttachmentAction   *SaveAttachment
	GetAttachmentAction    *GetAttachment
	DeleteAttachmentAction *DeleteAttachment
}

func (facade *AttachmentActionsFacade) SaveAttachment(ctx context.Context, attachment *Attachment) error {
	return facade.SaveAttachmentAction.Execute(ctx, attachment)
}

func (facade *AttachmentActionsFacade) GetAttachmentBy(ctx context.Context, attachmentId string) (*Attachment, error) {
	return facade.GetAttachmentAction.GetOneBy(ctx, attachmentId)
}

func (facade *AttachmentActionsFacade) GetAttachments(ctx context.Context, budgetId string, budgetType BudgetType) ([]AttachmentMetadata, error) {
	return facade.GetAttachmentAction.GetAllBy(ctx, budgetId, budgetType)
}

func (facade *AttachmentActionsFacade) DeleteAttachment(ctx context.Context, attachmentId string) error {
	return facade.DeleteAttachmentAction.Execute(ctx, attachmentId)
}
