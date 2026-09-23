//go:build test

package dynamodb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mrflick72/budget/budget-api/domain/budget/scheduledexpense"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/stretchr/testify/mock"
)

// client is shared across all tests in this package; each test uses its own
// unique stubbed user (see testutils.NewStubbedContextWith in each test) so
// FindAll's per-user partition scoping keeps tests isolated from one another
// on the shared table, instead of relying on cleanup between tests.
var client, _ = newDynamoDBClient()

type DynamoDbScheduledExpenseIdProviderMock struct {
	mock.Mock
}

func (mock *DynamoDbScheduledExpenseIdProviderMock) GenerateIdFor(se *scheduledexpense.ScheduledExpense) scheduledexpense.ScheduledExpenseId {
	args := mock.Called(se)
	return args.String(0)
}

func newDynamoDBClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("xxx", "xxx", "xxx")),
		config.WithRegion("eu-central-1"),
		config.WithBaseEndpoint("http://localhost:4566"),
	)

	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	return dynamodb.NewFromConfig(cfg), err
}

var TableName = "BUDGET_SCHEDULED_EXPENSE_TABLE_NAME_STAGING"

// newScheduledExpenseRepository wires a repository with an unconfigured tag
// mock — fine for untagged saves, whose empty "tag" attribute resolves to
// UNKNOWN without ever calling the tag repository.
func newScheduledExpenseRepository(idProvider scheduledexpense.ScheduledExpenseIdProvider) *DynamoDbScheduledExpenseRepository {
	return newScheduledExpenseRepositoryWith(idProvider, new(tags.SearchTagRepositoryMock))
}

// newScheduledExpenseRepositoryWith lets a test supply a configured tag mock
// to assert key→value resolution on read.
func newScheduledExpenseRepositoryWith(idProvider scheduledexpense.ScheduledExpenseIdProvider, searchTagRepository tags.SearchTagRepository) *DynamoDbScheduledExpenseRepository {
	return NewDynamoDbScheduledExpenseRepository(TableName, client, idProvider, searchTagRepository).(*DynamoDbScheduledExpenseRepository)
}

func setupTestDynamoDBTable() error {
	teardownTestDynamoDBTable()
	_, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(TableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("user_name"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("user_name"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeRange,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})

	if err != nil {
		var resourceInUseException *types.ResourceInUseException
		if !errors.As(err, &resourceInUseException) {
			return err
		}
	}
	return nil
}

func teardownTestDynamoDBTable() error {
	c, _ := newDynamoDBClient()
	_, err := c.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(TableName),
	})
	return err
}
