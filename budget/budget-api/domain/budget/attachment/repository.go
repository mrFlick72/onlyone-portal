package attachment

import "context"

type AttachmentRepository interface {
	SaveAttachment(ctx context.Context, attachment *Attachment) error
	GenAttachment(ctx context.Context, attachmentId string) (*Attachment, error)
	FindAllAttachment(ctx context.Context, budgetId string, budgetType BudgetType) ([]AttachmentMetadata, error)
	DeleteAttachment(ctx context.Context, attachmentId string) error
}
