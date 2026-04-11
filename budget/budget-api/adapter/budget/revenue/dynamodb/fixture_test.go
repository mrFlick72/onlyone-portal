package dynamodb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mrflick72/budget/budget-api/domain/budget/revenue"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/stretchr/testify/mock"
)

var client, _ = newDynamoDBClient()
var ctx = testutils.NewStubbedContextWith(security.UserName("testuser"))
var ctxAnotherUser = testutils.NewStubbedContextWith(security.UserName("anotheruser"))

type DynamoDbRevenueIdProviderMock struct {
	mock.Mock
}

func (mock *DynamoDbRevenueIdProviderMock) GenerateIdFor(revenue *revenue.Revenue) string {
	args := mock.Called(revenue)
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

var TableName = "BUDGET_REVENUE_TABLE_NAME_STAGING"

func newRevenueRepository(revenueIdProvider revenue.RevenueIdProvider) *DynamoDbRevenueRepository {
	return NewDynamoDbRevenueRepository(TableName, client, revenueIdProvider).(*DynamoDbRevenueRepository)
}

func setupTestDynamoDBTable() error {
	teardownTestDynamoDBTable()
	_, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(TableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("pk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("range_key"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("pk"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("range_key"),
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
