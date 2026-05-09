package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
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
	fileLocation string,
) error {
	item := map[string]types.AttributeValue{
		"pk":            &types.AttributeValueMemberS{Value: pk},
		"attachment_id": &types.AttributeValueMemberS{Value: att.AttachmentId},
		"budget_id":     &types.AttributeValueMemberS{Value: att.BudgetId},
		"budget_type":   &types.AttributeValueMemberS{Value: att.BudgetType},
		"date":          &types.AttributeValueMemberS{Value: att.Date.GetIsoFormattedDate()},
		"user_name":     &types.AttributeValueMemberS{Value: att.Owner},
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

func (repository *DynamoDbAttachmentMetadataRepository) FindAllAttachment(ctx context.Context, owner string, pk string) ([]attachment.AttachmentMetadata, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(repository.TableName),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":        &types.AttributeValueMemberS{Value: pk},
			":user_name": &types.AttributeValueMemberS{Value: owner},
		},
		KeyConditionExpression: aws.String("pk =:pk"),
		FilterExpression:       aws.String("user_name =:user_name"),
	}

	query, err := repository.Client.Query(ctx, input)
	if err != nil {
		repository.logger.LogErrorfFor("Error during DynamoDb Query: $v", err)
		return nil, err
	}

	result := make([]attachment.AttachmentMetadata, 0, len(query.Items))
	for _, item := range query.Items {
		rawDate := item["date"].(*types.AttributeValueMemberS).Value
		dateFor, err := date.IsoDateFor(rawDate)
		if err != nil {
			errorMessage := fmt.Sprint("Parsing failure: the data field is in a incorrect format: raw data: %s Error details: %v", rawDate, err)
			repository.logger.LogErrorfFor(errorMessage)
			return nil, errors.New(errorMessage)
		}
		result = append(result, attachment.AttachmentMetadata{
			AttachmentId: item["attachment_id"].(*types.AttributeValueMemberS).Value,
			BudgetId:     item["budget_id"].(*types.AttributeValueMemberS).Value,
			BudgetType:   item["budget_type"].(*types.AttributeValueMemberS).Value,
			Date:         *dateFor,
			Owner:        item["user_name"].(*types.AttributeValueMemberS).Value,
			FineName:     item["file_name"].(*types.AttributeValueMemberS).Value,
			ContentType:  item["content_type"].(*types.AttributeValueMemberS).Value,
			Metadata:     map[string]string{},
		})
	}

	return result, nil
}

// GetAttachment looks up an attachment by id via the Global Secondary Index.
// The GSI does not project budget_id, budget_type or date, so the returned
// Attachment is intentionally incomplete: only AttachmentId, Owner, FineName
// and ContentType are populated. Callers that need the full metadata must
// fetch it from the base table.
func (repository *DynamoDbAttachmentMetadataRepository) GetAttachment(ctx context.Context, user string, id string) (*attachment.Attachment, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(repository.TableName),
		IndexName: aws.String(fmt.Sprintf("%s_GLOBAL_INDEX", repository.TableName)),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":attachmentId": &types.AttributeValueMemberS{Value: id},
			":user_name":    &types.AttributeValueMemberS{Value: user},
		},
		KeyConditionExpression: aws.String("attachmentId =:attachmentId"),
		FilterExpression:       aws.String("user_name =:user_name"),
	}
	query, err := repository.Client.Query(ctx, input)
	if err != nil {
		repository.logger.LogErrorfFor("Error during DynamoDb Global Secondary Index query err: %v", err)
		return nil, err
	}

	if len(query.Items) == 0 {
		repository.logger.LogErrorfFor("No attachment found in the global secondary index with this id: %s", id)
		return nil, errors.New("attachment not found")
	}

	attachmentMetadata := query.Items[0]
	return &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			AttachmentId: id,
			Owner:        user,
			FineName:     attachmentMetadata["file_name"].(*types.AttributeValueMemberS).Value,
			ContentType:  attachmentMetadata["content_type"].(*types.AttributeValueMemberS).Value,
		},
	}, nil
}
