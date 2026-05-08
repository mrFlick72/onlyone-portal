package attachment

import (
	"context"
	"errors"

	"github.com/mrflick72/budget/budget-api/adapter/budget/attachment/dynamodb"
	"github.com/mrflick72/budget/budget-api/adapter/budget/attachment/s3"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
)

type AwsCompositeAttachmentRepository struct {
	IdProvider         *dynamodb.DynamoDbAttachmentIdProvider
	MetadataRepository *dynamodb.DynamoDbAttachmentMetadataRepository
	ContentRepository  *s3.S3AttachmentContentRepository
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
	}
}

func (repository *AwsCompositeAttachmentRepository) SaveAttachment(ctx context.Context, att *attachment.Attachment) error {
	att.AttachmentId = repository.IdProvider.GenerateAttachmentIdFor(att)

	pk := repository.IdProvider.PartitionKeyFor(att)
	rk := repository.IdProvider.RangeKeyFor(att.AttachmentId)

	objectKey := repository.ContentRepository.ObjectKeyFor(att, pk, rk)
	fileLocation := repository.ContentRepository.FileLocationFor(objectKey)

	if err := repository.ContentRepository.Save(ctx, objectKey, att.ContentType, att.Content); err != nil {
		return err
	}

	return repository.MetadataRepository.Save(ctx, att, pk, rk, fileLocation)
}

func (repository *AwsCompositeAttachmentRepository) GenAttachment(ctx context.Context, attachmentId string) (*attachment.Attachment, error) {
	return nil, errors.New("not implemented")
}

func (repository *AwsCompositeAttachmentRepository) FindAllAttachment(ctx context.Context, budgetId string, budgetType attachment.BudgetType) ([]attachment.AttachmentMetadata, error) {
	repository.MetadataRepository.FindAllAttachment(ctx, budgetType)
	return nil, errors.New("not implemented")
}

func (repository *AwsCompositeAttachmentRepository) DeleteAttachment(ctx context.Context, attachmentId string) error {
	return errors.New("not implemented")
}
