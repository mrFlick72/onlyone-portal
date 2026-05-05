package attachment

import "context"

type AttachmentRepository interface {
	SaveAttachment(ctx context.Context, attachment *Attachment) error
	GenAttachment(ctx context.Context, attachmentId string) (error, *Attachment)
	FindAllAttachment(ctx context.Context, budgetId string) (error, []AttachmentMetadata)
	DeleteAttachment(ctx context.Context, attachmentId string) error
}
