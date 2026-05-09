//go:build test

package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mrflick72/budget/budget-api/internal/testutils/awsfixture"
)

var TableName = "BUDGET_ATTACHMENT_METADATA_TABLE_NAME_STAGING"

var client *dynamodb.Client = awsfixture.NewLocalDynamoDBClient()

func setupTestDynamoDBTable() error {
	_ = awsfixture.TeardownTable(context.TODO(), client, TableName)
	return awsfixture.SetupAttachmentTable(context.TODO(), client, TableName)
}

func teardownTestDynamoDBTable() error {
	return awsfixture.TeardownTable(context.TODO(), client, TableName)
}

func getItem(ctx context.Context, pk, attachmentId string) (map[string]types.AttributeValue, error) {
	out, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"pk":            &types.AttributeValueMemberS{Value: pk},
			"attachment_id": &types.AttributeValueMemberS{Value: attachmentId},
		},
	})
	if err != nil {
		return nil, err
	}
	return out.Item, nil
}
