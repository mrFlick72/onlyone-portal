package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

// @Test
// public void saveAnewBudgetExpense() {
//     BudgetExpenseId id = new BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU");
//     BudgetExpense expected = new BudgetExpense(id, new UserName("USER"), DATE, Money.moneyFor("10.50"), "NOTE", "TAG");

//     BudgetExpense actual = budgetExpenseRepository.save(expected);

//     Assertions.assertEquals(expected, actual);
//     BudgetExpense retrievedBudgetExpense = budgetExpenseRepository.findFor(expected.id()).get();
//     Assertions.assertEquals(expected, retrievedBudgetExpense);
// }

// @Test
// public void findByDateRange() {
//     budgetExpenseRepository.save(new BudgetExpense(null, new UserName("USER"), Date.dateFor("06/05/2018"), Money.moneyFor("17.50"), "Lanch", "lanch"));
//     List<BudgetExpense> actualRange = budgetExpenseRepository.findByDateRange(new UserName("USER"), Date.dateFor("01/01/2018"), Date.dateFor("05/05/2018"));
//     List<BudgetExpense> expectedRange =
//             asList(
//                     new BudgetExpense(null, new UserName("USER"), Date.dateFor("06/01/2018"), Money.moneyFor("17.50"), "Lanch", "lanch"),
//                     new BudgetExpense(null, new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("10.50"), "Super Market", "super-market"),
//                     new BudgetExpense(null, new UserName("USER"), Date.dateFor("13/02/2018"), Money.moneyFor("17.50"), "Dinner", "dinner"),
//                     new BudgetExpense(null, new UserName("USER"), Date.dateFor("22/02/2018"), Money.moneyFor("17.50"), "Super Market", "super-market"),
//                     new BudgetExpense(null, new UserName("USER"), Date.dateFor("05/05/2018"), Money.moneyFor("17.50"), "Lanch", "lanch")
//             );

//     Assertions.assertEquals(expectedRange.size(), actualRange.size());
// }

// @Test
// public void findByDateRangeAndSearchTags() {
//     List<BudgetExpense> actualRange = budgetExpenseRepository.findByDateRange(new UserName("USER"), Date.dateFor("01/01/2018"), Date.dateFor("05/05/2018"), "super-market", "dinner");
//     List<BudgetExpense> expectedRange =
//             asList(new BudgetExpense(new BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU"), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("10.50"), "Super Market", "super-market"),
//                     new BudgetExpense(new BudgetExpenseId("MjAxOF8yX1VTRVI=-MTNfQV9TQUxU"), new UserName("USER"), Date.dateFor("13/02/2018"), Money.moneyFor("17.50"), "Dinner", "dinner"),
//                     new BudgetExpense(new BudgetExpenseId("MjAxOF8yX1VTRVI=-MjJfQV9TQUxU"), new UserName("USER"), Date.dateFor("22/02/2018"), Money.moneyFor("17.50"), "Super Market", "super-market"));

//     Assertions.assertTrue(expectedRange.containsAll(actualRange) );
// }
// @Test
// public void deleteBudgetExpense() {
//     BudgetExpenseId id = new BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU");
//     BudgetExpense expected = new BudgetExpense(id, new UserName("USER"), DATE, Money.moneyFor("10.50"), "NOTE", "TAG");

//     budgetExpenseRepository.save(expected);

//     budgetExpenseRepository.delete(id);
//     Optional<BudgetExpense> actual = budgetExpenseRepository.findFor(id);

//     Assertions.assertThrows(Exception.class, actual::orElseThrow);
// }

var TableName = "BUDGET_EXPENSE_TABLE_NAME_STAGING"
var client, _ = newDynamoDBClient()

func newStubbedContext() *context.Context {
	ctx := context.Background()
	userName := security.UserName("testuser")
	user := security.User{UserName: &userName}
	newCtx := context.WithValue(ctx, "user", user)

	return &newCtx
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

func TestMain(m *testing.M) {
	setupTestDynamoDBTable()

	code := m.Run() // run all tests

	teardownTestDynamoDBTable()
	os.Exit(code)
}

// @Test
// public void saveAnewBudgetExpense() {
//     BudgetExpenseId id = new BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU");
//     BudgetExpense expected = new BudgetExpense(id, new UserName("USER"), DATE, Money.moneyFor("10.50"), "NOTE", "TAG");

//     BudgetExpense actual = budgetExpenseRepository.save(expected);

//     Assertions.assertEquals(expected, actual);
//     BudgetExpense retrievedBudgetExpense = budgetExpenseRepository.findFor(expected.id()).get();
//     Assertions.assertEquals(expected, retrievedBudgetExpense);
// }

func TestSaveANewBudgetExpense(t *testing.T) {
	mockedBudgetExpenseIdProvider := new(DynamoDbBudgetExpenseIdProviderMock)
	repo := newBudgetExpenseRepository(mockedBudgetExpenseIdProvider)
	ctx := newStubbedContext()

	// Implement the test logic here
	expected := expense.BudgetExpense{
		Id:       expense.BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU"),
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tag:      "TAG",
	}

	mockedBudgetExpenseIdProvider.On("GenerateIdFor", &expected).Return("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU")

	err := repo.Save(ctx, &expected)

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	retrievedBudgetExpense, err := repo.FindFor(ctx, expected.Id)
	if err != nil {
		t.Errorf("Error retrieving budget expense: %v", err)
	}
	fmt.Println("Expected:", expected.Amount.StringifyAmount())
	fmt.Println("Retrieved:", retrievedBudgetExpense.Amount.StringifyAmount())
	assert.Equal(t, expected, retrievedBudgetExpense)
}
