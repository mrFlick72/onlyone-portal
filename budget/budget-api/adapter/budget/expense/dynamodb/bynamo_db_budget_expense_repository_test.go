package dynamodb

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestMain(m *testing.M) {
	setupTestDynamoDBTable()

	code := m.Run() // run all tests

	teardownTestDynamoDBTable()
	os.Exit(code)
}

func TestFindBudgetExpenseByDateRange(t *testing.T) {	
	mockedBudgetExpenseIdProvider := &DynamoDbBudgetExpenseIdProvider{
		SaltGenerator: func() string { return uuid.New().String() },
	}
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	searchTagMockRepositorySetup(ctxUser, mockSearchTagRepository)

	repo := newBudgetExpenseRepository(mockedBudgetExpenseIdProvider, mockSearchTagRepository)
	err := loadBudgetExpensesFromCSVFile("find-by-date-range-data-set.csv", mockedBudgetExpenseIdProvider, mockSearchTagRepository)
	if err != nil {
		t.Fatalf("Failed to load budget expenses: %v", err)
	}

	result, err := repo.FindByDateRange(ctxUser, testutils.SafeDateFor("01/02/2018"), testutils.SafeDateFor("28/02/2019"), []tags.SearchTagKey{})
	assert.Equal(t, nil, err)
	assert.Equal(t, 11, len(result))
}



func TestFindNonExistentBudgetExpenseReturnsNil(t *testing.T) {
	mockedBudgetExpenseIdProvider := new(DynamoDbBudgetExpenseIdProviderMock)
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	repo := newBudgetExpenseRepository(mockedBudgetExpenseIdProvider, mockSearchTagRepository)

	// Implement the test logic here
	nonExistentBudgetExpenseId := expense.BudgetExpenseId("dyghtq4hrbg-MTJfQV9TQUxU")

	result, err := repo.FindFor(ctx, nonExistentBudgetExpenseId)

	assert.NotEqual(t, nil, err)
	assert.Equal(t, nil, result)
}

func TestFindBudgetExpenseOfOtherPersonISNotAllowed(t *testing.T) {
	mockedBudgetExpenseIdProvider := new(DynamoDbBudgetExpenseIdProviderMock)
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	repo := newBudgetExpenseRepository(mockedBudgetExpenseIdProvider, mockSearchTagRepository)

	// Implement the test logic here
	input := expense.BudgetExpense{
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tag:      []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
	}

	// Implement the test logic here
	expected := expense.BudgetExpense{
		Id:       expense.BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU"),
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tag:      []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
	}

	mockedBudgetExpenseIdProvider.On("GenerateIdFor", &input).Return("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU")

	err := repo.Save(testutils.NewStubbedContextWith("another User"), &input)

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	retrievedBudgetExpense, err := repo.FindFor(ctxAnotherUser, expected.Id)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, nil, retrievedBudgetExpense)
}

func TestSaveANewBudgetExpense(t *testing.T) {
	mockedBudgetExpenseIdProvider := new(DynamoDbBudgetExpenseIdProviderMock)
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	repo := newBudgetExpenseRepository(mockedBudgetExpenseIdProvider, mockSearchTagRepository)
	searchTagMockRepositorySetup(ctx, mockSearchTagRepository)

	// Implement the test logic here
	input := expense.BudgetExpense{
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tag:      []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
	}

	// Implement the test logic here
	expected := expense.BudgetExpense{
		Id:       expense.BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU"),
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tag:      []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
	}

	mockedBudgetExpenseIdProvider.On("GenerateIdFor", &input).Return("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU")

	err := repo.Save(ctx, &input)

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	retrievedBudgetExpense, err := repo.FindFor(ctx, expected.Id)
	if err != nil {
		t.Errorf("Error retrieving budget expense: %v", err)
	}

	mockedBudgetExpenseIdProvider.AssertCalled(t, "GenerateIdFor", &input)
	assert.Equal(t, expected, retrievedBudgetExpense)
}

func TestUpdateABudgetExpense(t *testing.T) {
	mockedBudgetExpenseIdProvider := new(DynamoDbBudgetExpenseIdProviderMock)
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	repo := newBudgetExpenseRepository(mockedBudgetExpenseIdProvider, mockSearchTagRepository)
	searchTagMockRepositorySetup(ctx, mockSearchTagRepository)

	//	Implement the test logic here
	input := expense.BudgetExpense{
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tag:      []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
	}
	mockedBudgetExpenseIdProvider.On("GenerateIdFor", &input).Return("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU")

	repo.Save(ctx, &input)

	// Implement the test logic here
	expected := expense.BudgetExpense{
		Id:       expense.BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU"),
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tag:      []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
	}

	err := repo.Save(ctx, &expected)

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	retrievedBudgetExpense, err := repo.FindFor(ctx, expected.Id)
	if err != nil {
		t.Errorf("Error retrieving budget expense: %v", err)
	}

	assert.Equal(t, expected, retrievedBudgetExpense)
	mockedBudgetExpenseIdProvider.AssertNumberOfCalls(t, "GenerateIdFor", 1)
}

func TestUpdateABudgetExpenseFailsWhenTheBudgetExpenseDoesNotBelongsToTheUserInTheContext(t *testing.T) {
	mockedBudgetExpenseIdProvider := new(DynamoDbBudgetExpenseIdProviderMock)
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	repo := newBudgetExpenseRepository(mockedBudgetExpenseIdProvider, mockSearchTagRepository)
	searchTagMockRepositorySetup(ctx, mockSearchTagRepository)

	//	Implement the test logic here
	input := expense.BudgetExpense{
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tag:      []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
	}
	mockedBudgetExpenseIdProvider.On("GenerateIdFor", &input).Return("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU")

	repo.Save(ctx, &input)

	// Implement the test logic here
	expected := expense.BudgetExpense{
		Id:       expense.BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU"),
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tag:      []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
	}

	err := repo.Save(ctxAnotherUser, &expected)

	mockedBudgetExpenseIdProvider.AssertNotCalled(t, "GenerateIdFor", &expected)
	fmt.Print(err)
	assert.NotEqual(t, nil, err)
}

func TestDeleteBudgetExpense(t *testing.T) {
	mockedBudgetExpenseIdProvider := new(DynamoDbBudgetExpenseIdProviderMock)
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	repo := newBudgetExpenseRepository(mockedBudgetExpenseIdProvider, mockSearchTagRepository)
	searchTagMockRepositorySetup(ctx, mockSearchTagRepository)

	// Implement the test logic here
	budgetExpense := expense.BudgetExpense{
		Id:       expense.BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU"),
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tag:      []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
	}
	mockedBudgetExpenseIdProvider.On("GenerateIdFor", &budgetExpense).Return("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU")

	err := repo.Save(ctx, &budgetExpense)

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	err = repo.Delete(ctx, budgetExpense.Id)
	if err != nil {
		t.Errorf("Expected nil error on delete, got %v", err)
	}

	_, err = repo.FindFor(ctx, budgetExpense.Id)
	if err == nil {
		t.Errorf("Expected error retrieving deleted budget expense, got nil")
	}
}

func TestDeleteBudgetExpenseFailsWhenTheBudgetExpenseDoesNotBelongsToTheUserInTheContext(t *testing.T) {
	mockedBudgetExpenseIdProvider := new(DynamoDbBudgetExpenseIdProviderMock)
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	repo := newBudgetExpenseRepository(mockedBudgetExpenseIdProvider, mockSearchTagRepository)
	searchTagMockRepositorySetup(ctxAnotherUser, mockSearchTagRepository)

	// Implement the test logic here
	testUserBudgetExpense := expense.BudgetExpense{
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tag:      []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
	}
	mockedBudgetExpenseIdProvider.On("GenerateIdFor", &testUserBudgetExpense).Return("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU")

	_ = repo.Save(ctx, &testUserBudgetExpense)

	// Implement the test logic here
	anotherUserBudgetExpense := expense.BudgetExpense{
		UserName: "anotheruser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tag:      []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
	}
	mockedBudgetExpenseIdProvider.On("GenerateIdFor", &anotherUserBudgetExpense).Return("MjAxOF8yX1VTRVI=-MjAxOF8yX1VTRVI")
	_ = repo.Save(ctxAnotherUser, &anotherUserBudgetExpense)

	// the user in the context is "testuser", but we try to delete another user's budget expense
	err := repo.Delete(ctx, anotherUserBudgetExpense.Id)
	assert.NotEqual(t, nil, err)

	// verify that the another user's budget expense is still there
	expected, err := repo.FindFor(ctxAnotherUser, anotherUserBudgetExpense.Id)
	assert.Equal(t, nil, err)
	assert.Equal(t, anotherUserBudgetExpense, expected)
}
