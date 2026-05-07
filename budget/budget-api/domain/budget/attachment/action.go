package attachment

import (
	"context"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

var logger = logging.GetLoggerInstanceForComponentByTypeName("attachment.SaveAttachment")

type SaveAttachment struct {
	repository AttachmentRepository
}

func NewSaveAttachment(repository AttachmentRepository) *SaveAttachment {
	return &SaveAttachment{repository: repository}
}

func (a *SaveAttachment) Execute(ctx context.Context, attachment *Attachment) error {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		logger.LogErrorfFor("error getting current user: %v\n", err)
		return err
	}
	attachment.Owner = *user.UserName
	if err := a.repository.SaveAttachment(ctx, attachment); err != nil {
		logger.LogErrorfFor("error saving attachment: %v\n", err)
		return err
	}
	return nil
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

func (a *GetAttachment) GetOneBy(ctx context.Context, attachmentId string) (*Attachment, error) {
	panic("has to be implemented")
}

func (a *GetAttachment) GetAllBy(ctx context.Context, budgetId string) ([]AttachmentMetadata, error) {
	panic("has to be implemented")
}
