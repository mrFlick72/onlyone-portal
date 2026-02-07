package dynamodb

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/stretchr/testify/mock"
)

var client, _ = newDynamoDBClient()
var ctx = testutils.NewStubbedContextWith(security.UserName("testuser"))
var ctxAnotherUser = testutils.NewStubbedContextWith(security.UserName("anotheruser"))

type DynamoDbBudgetExpenseIdProviderMock struct {
	mock.Mock
}

func (mock *DynamoDbBudgetExpenseIdProviderMock) GenerateIdFor(budgetExpense *expense.BudgetExpense) string {
	args := mock.Called(budgetExpense)
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

var TableName = "BUDGET_EXPENSE_TABLE_NAME_STAGING"

func newBudgetExpenseRepository(budgetExpenseIdProvider expense.BudgetExpenseIdProvider) *DynamoDbBudgetExpenseRepository {
	return &DynamoDbBudgetExpenseRepository{
		// Initialize with mock or test dependencies as needed
		TableName:               TableName,
		Client:                  client,
		BudgetExpenseIdProvider: budgetExpenseIdProvider,
	}
}
func setupTestDynamoDBTable() error {
	// it is an attempt to clean up possible dirty state before creating
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
	client, _ := newDynamoDBClient()
	_, err := client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(TableName),
	})
	return err
}

func loadBudgetExpensesFromCSVFile(filePath string, mockedBudgetExpenseIdProvider *DynamoDbBudgetExpenseIdProvider) error {
	recordsNumber := 0
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return err
	}

	repository := newBudgetExpenseRepository(mockedBudgetExpenseIdProvider)

	for _, record := range records {
		recordsNumber++

		budgetExpense := expense.BudgetExpense{
			UserName: security.UserName(record[0]),
			Date:     testutils.SafeDateFor(record[1]),
			Amount:   testutils.SafeMoneyFor(record[2]),
			Note:     record[3],
			Tag:      record[4],
		}

		err := repository.Save(testutils.NewStubbedContextWith(record[0]), &budgetExpense)
		fmt.Printf("Loaded budget expense id: %+v\n", budgetExpense.Id)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Loaded %d budget expenses from %s\n", recordsNumber, filePath)
	return nil
}