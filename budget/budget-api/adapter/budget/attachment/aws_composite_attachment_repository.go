package attachment

import (
	"context"
	"errors"

	"github.com/mrflick72/budget/budget-api/adapter/budget/attachment/dynamodb"
	"github.com/mrflick72/budget/budget-api/adapter/budget/attachment/s3"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type AwsCompositeAttachmentRepository struct {
	IdProvider         *dynamodb.DynamoDbAttachmentIdProvider
	MetadataRepository *dynamodb.DynamoDbAttachmentMetadataRepository
	ContentRepository  *s3.S3AttachmentContentRepository
	logger             *logging.Logger
}

func NewAwsCompositeAttachmentRepository(
	idProvider *dynamodb.DynamoDbAttachmentIdProvider,
	metadataRepository *dynamodb.DynamoDbAttachmentMetadataRepository,
	contentRepository *s3.S3AttachmentContentRepository,
) attachment.AttachmentRepository {
	return &AwsCompositeAttachmentRepository{
		IdProvider:         idProvider,
		MetadataRepository: metadataRepository,
		ContentRepository:  contentRepository,
		logger:             logging.GetLoggerInstanceForComponentByType(&AwsCompositeAttachmentRepository{}),
	}
}

func (repository *AwsCompositeAttachmentRepository) SaveAttachment(ctx context.Context, att *attachment.Attachment) error {
	att.AttachmentId = repository.IdProvider.GenerateAttachmentIdFor(att)

	pk := repository.IdProvider.PartitionKeyFor(att)

	objectKey := repository.ContentRepository.ObjectKeyFor(att, pk)
	fileLocation := repository.ContentRepository.FileLocationFor(objectKey)

	if err := repository.ContentRepository.Save(ctx, objectKey, att.ContentType, att.Content); err != nil {
		return err
	}

	return repository.MetadataRepository.Save(ctx, att, pk, fileLocation)
}

func (repository *AwsCompositeAttachmentRepository) GenAttachment(ctx context.Context, attachmentId string) (*attachment.Attachment, error) {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		repository.logger.LogErrorfFor("error getting current user: %v", err)
		return nil, err
	}
	att, err := repository.MetadataRepository.GetAttachment(ctx, *user.UserName, attachmentId)
	repository.logger.LogInfofFor("attachment metadata: %v", att)
	if err != nil {
		repository.logger.LogErrorfFor("error getting attachment metadata for id %s: %v", attachmentId, err)
		return nil, err
	}

	pk := repository.IdProvider.PartitionKeyFor(att)
	objectKey := repository.ContentRepository.ObjectKeyFor(att, pk)
	repository.logger.LogInfofFor("objectKey: %v", objectKey)

	content, err := repository.ContentRepository.GetContentFor(ctx, objectKey)
	repository.logger.LogInfofFor("content size: %v", len(content))

	if err != nil {
		repository.logger.LogErrorfFor("error getting attachment content for id %s: %v", attachmentId, err)
		return nil, err
	}
	att.Content = content
	return att, nil
}

func (repository *AwsCompositeAttachmentRepository) FindAllAttachment(ctx context.Context, budgetId string, budgetType attachment.BudgetType) ([]attachment.AttachmentMetadata, error) {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		repository.logger.LogErrorfFor("error getting current user: %v", err)
		return []attachment.AttachmentMetadata{}, err
	}
	att := attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			BudgetId:   budgetId,
			BudgetType: string(budgetType),
		},
	}
	pk := repository.IdProvider.PartitionKeyFor(&att)
	repository.logger.LogDebugfFor("querying attachments for partition key %s", pk)
	return repository.MetadataRepository.FindAllAttachment(ctx, *user.UserName, pk)
}

func (repository *AwsCompositeAttachmentRepository) DeleteAttachment(ctx context.Context, attachmentId string) error {
	return errors.New("not implemented")
}
