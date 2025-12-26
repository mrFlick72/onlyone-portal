package test 

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var client, _ = newDynamoDBClient()
var TableName = "Tags_Local_Test_Table"

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

func setupTestDynamoDBTable() error {
	deleteTestDynamoDBTable()
	createTestDynamoDBTable()
	return nil
}

func deleteTestDynamoDBTable() error {
	_, err := client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(TableName),
	})
	if err != nil {
		var resourceNotFoundException *types.ResourceNotFoundException
		fmt.Println("Error deleting table:", err)
		if !errors.As(err, &resourceNotFoundException) {
			return err
		}
	}
	return nil
}

func createTestDynamoDBTable() error {
	_, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(TableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("user_name"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("search_tag_key"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("user_name"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("search_tag_key"),
				KeyType:       types.KeyTypeRange,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		var resourceInUseException *types.ResourceInUseException
		fmt.Println("Error creating table:", err)
		if !errors.As(err, &resourceInUseException) {
			return err
		}
	}
	return nil
}

func main() {
	setupTestDynamoDBTable()
}
