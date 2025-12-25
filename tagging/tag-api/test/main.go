package main

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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
	// Create table if not exists
	client, _ := newDynamoDBClient()

	TableName := "Tags_Local_Test_Table"
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

func main() {
	setupTestDynamoDBTable()
}
