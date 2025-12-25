package dynamodb

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/domain"
)

var TableName = "TestTagsTable"

func newStubbedContext() *context.Context {
	ctx := context.Background()
	userName := domain.UserName("testuser")
	user := domain.User{UserName: &userName}
	newCtx, _ := domain.SetCurrentUser(user, &ctx)
	return newCtx
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

func newTagDynamoDBRepository() *TagDynamoDBRepository {
	client, _ := newDynamoDBClient()
	return &TagDynamoDBRepository{
		// Initialize with mock or test dependencies as needed
		TableName: TableName,
		Client:    client,
	}
}

func setupTestDynamoDBTable() error {
	// Create table if not exists
	client, _ := newDynamoDBClient()

	_, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(TableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("UserName"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("Key"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("UserName"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("Key"),
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
	client, _ := newDynamoDBClient()
	_, err := client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(TableName),
	})
	return err
}

func TestMain(m *testing.M) {
	setupTestDynamoDBTable()

	code := m.Run() // run all tests

	teardownTestDynamoDBTable()
	os.Exit(code)
}

func saveATag(t *testing.T) {
	repo := newTagDynamoDBRepository()
	tag := domain.Tag{Key: "exampleKey", Value: "exampleValue"}

	err := repo.SaveTag(newStubbedContext(), &tag)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSaveTag(t *testing.T) {
	saveATag(t)
	repo := newTagDynamoDBRepository()
	key := "exampleKey"

	tag, err := repo.GetTagBy(newStubbedContext(), key)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if tag.Key != key {
		t.Errorf("Expected key %v, got %v", key, tag.Key)
	}
}

func TestFindAllTags(t *testing.T) {
	saveATag(t)

	repo := newTagDynamoDBRepository()

	tags, err := repo.FindAllTags(newStubbedContext())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(*tags) == 0 {
		t.Errorf("Expected some tags, got none")
	}
}
