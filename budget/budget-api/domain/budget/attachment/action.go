package attachment

import "context"

type SaveAttachment struct {
	repository AttachmentRepository
}

func (a *SaveAttachment) Execute(ctx context.Context, attachment *Attachment) error {
	panic("has to be implemented")
}

type DeleteAttachment struct {
	repository AttachmentRepository
}

func (a *DeleteAttachment) Execute(ctx context.Context, attachmentId string) error {
	panic("has to be implemented")
}

type GetAttachment struct {
	repository AttachmentRepository
}

func (a *GetAttachment) GetOneBy(ctx context.Context, attachmentId string) (error, *Attachment) {
	panic("has to be implemented")
}

func (a *GetAttachment) GetAllBy(ctx context.Context, budgetId string) (error, []AttachmentMetadata) {
	panic("has to be implemented")
}
