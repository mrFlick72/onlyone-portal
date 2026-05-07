package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
)

type DynamoDbAttachmentMetadataRepository struct {
	TableName  string
	Client     *dynamodb.Client
	IdProvider *DynamoDbAttachmentIdProvider
	logger     *logging.Logger
}

func NewDynamoDbAttachmentMetadataRepository(
	tableName string,
	client *dynamodb.Client,
	idProvider *DynamoDbAttachmentIdProvider,
) *DynamoDbAttachmentMetadataRepository {
	return &DynamoDbAttachmentMetadataRepository{
		TableName:  tableName,
		Client:     client,
		IdProvider: idProvider,
		logger:     logging.GetLoggerInstanceForComponentByType(&DynamoDbAttachmentMetadataRepository{}),
	}
}

func (repository *DynamoDbAttachmentMetadataRepository) Save(
	ctx context.Context,
	att *attachment.Attachment,
	pk string,
	rk string,
	fileLocation string,
) error {
	item := map[string]types.AttributeValue{
		"pk":            &types.AttributeValueMemberS{Value: pk},
		"range_key":     &types.AttributeValueMemberS{Value: rk},
		"attachment_id": &types.AttributeValueMemberS{Value: att.AttachmentId},
		"budget_id":     &types.AttributeValueMemberS{Value: att.BudgetId},
		"budget_type":   &types.AttributeValueMemberS{Value: att.BudgetType},
		"date":          &types.AttributeValueMemberS{Value: att.Date.GetIsoFormattedDate()},
		"owner":         &types.AttributeValueMemberS{Value: att.Owner},
		"file_name":     &types.AttributeValueMemberS{Value: att.FineName},
		"content_type":  &types.AttributeValueMemberS{Value: att.ContentType},
		"file_location": &types.AttributeValueMemberS{Value: fileLocation},
	}

	if len(att.Metadata) > 0 {
		metadataMap := make(map[string]types.AttributeValue, len(att.Metadata))
		for k, v := range att.Metadata {
			metadataMap[k] = &types.AttributeValueMemberS{Value: v}
		}
		item["metadata"] = &types.AttributeValueMemberM{Value: metadataMap}
	}

	_, err := repository.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(repository.TableName),
		Item:      item,
	})
	if err != nil {
		repository.logger.LogErrorfFor("error saving attachment metadata for id %s: %v", att.AttachmentId, err)
	}
	return err
}
