package dynamodb

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

	result, err = repo.FindByDateRange(ctxUser, testutils.SafeDateFor("01/02/2018"), testutils.SafeDateFor("28/02/2019"), []tags.SearchTagKey{"dinner"})
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(result))

	result, err = repo.FindByDateRange(ctxUser, testutils.SafeDateFor("01/02/2018"), testutils.SafeDateFor("28/02/2019"), []tags.SearchTagKey{"dinner", "lunch"})
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(result))

	result, err = repo.FindByDateRange(ctxUser, testutils.SafeDateFor("01/02/2018"), testutils.SafeDateFor("28/02/2019"), []tags.SearchTagKey{"unknown-tag"})
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(result))
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
		Tags:     []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
	}

	// Implement the test logic here
	expected := expense.BudgetExpense{
		Id:       expense.BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU"),
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tags:     []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
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
		Tags:     []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
	}

	// Implement the test logic here
	expected := expense.BudgetExpense{
		Id:       expense.BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU"),
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tags:     []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
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
		Tags:     []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
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
		Tags:     []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
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
		Tags:     []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
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
		Tags:     []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
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
		Tags:     []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
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
		Tags:     []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
	}
	mockedBudgetExpenseIdProvider.On("GenerateIdFor", &testUserBudgetExpense).Return("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU")

	_ = repo.Save(ctx, &testUserBudgetExpense)

	// Implement the test logic here
	anotherUserBudgetExpense := expense.BudgetExpense{
		UserName: "anotheruser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tags:     []tags.SearchTag{{Key: "TAG", Value: "TAG"}},
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

// The detection tests below exercise fromDynamo directly with a hand-built item
// so they do not need a LocalStack round-trip and can hit every branch of the
// "a tag was deleted" condition deterministically. Each test injects its own
// bus so the emitted event can be asserted in isolation. A read distinguishes a
// deleted tag (a stored key that now resolves to UNKNOWN) from a record stored
// with the UNKNOWN sentinel as its key (the default applied at create time),
// which is not a deletion.

func TestWhenASearchTagIsBeenRemoved(t *testing.T) {
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockSearchTagRepository.On("GetTagBy", ctx, "deleted-tag").Return(searchTagPtr(tags.UnknownSentinel()), nil)

	bus := expense.NewEventBus()
	repo := newDetectionRepository(mockSearchTagRepository, bus)

	result, err := repo.fromDynamo(ctx, dynamoItemFor("deleted-tag"))

	assert.Equal(t, nil, err)
	assert.Equal(t, []tags.SearchTag{tags.UnknownSentinel()}, result.Tags)
	event := awaitReclassifyEvent(t, bus)
	assert.Equal(t, expense.BudgetExpenseId("A_BUDGET_ID"), event.Payload.Id)
	assert.Equal(t, []tags.SearchTag{tags.UnknownSentinel()}, event.Payload.Tags)
}

func TestWhenEverySearchTagWasRemovedTheReclassifiedPayloadIsDeduplicated(t *testing.T) {
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockSearchTagRepository.On("GetTagBy", ctx, "gone-a").Return(searchTagPtr(tags.UnknownSentinel()), nil)
	mockSearchTagRepository.On("GetTagBy", ctx, "gone-b").Return(searchTagPtr(tags.UnknownSentinel()), nil)

	bus := expense.NewEventBus()
	repo := newDetectionRepository(mockSearchTagRepository, bus)

	result, err := repo.fromDynamo(ctx, dynamoItemFor("gone-a,gone-b"))

	assert.Equal(t, nil, err)
	// The read output preserves one resolved tag per stored key...
	assert.Equal(t, []tags.SearchTag{tags.UnknownSentinel(), tags.UnknownSentinel()}, result.Tags)
	// ...while the durable-reclassify payload collapses the repeated sentinels.
	event := awaitReclassifyEvent(t, bus)
	assert.Equal(t, []tags.SearchTag{tags.UnknownSentinel()}, event.Payload.Tags)
}

func TestWhenNoSearchTagWasRemovedNoReclassifyEventIsFired(t *testing.T) {
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockSearchTagRepository.On("GetTagBy", ctx, "super-market").Return(&tags.SearchTag{Key: "super-market", Value: "super-market"}, nil)

	bus := expense.NewEventBus()
	repo := newDetectionRepository(mockSearchTagRepository, bus)

	result, err := repo.fromDynamo(ctx, dynamoItemFor("super-market"))

	assert.Equal(t, nil, err)
	assert.Equal(t, []tags.SearchTag{{Key: "super-market", Value: "super-market"}}, result.Tags)
	assertNoReclassifyEventFired(t, bus)
}

// A record legitimately stored with the UNKNOWN sentinel as its tag key (the
// default applied at create time) resolves to UNKNOWN but was not deleted, so
// it must not be reclassified — otherwise every untagged record would re-Save
// and re-emit on every read.
func TestWhenTheStoredTagIsTheUnknownSentinelNoReclassifyEventIsFired(t *testing.T) {
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockSearchTagRepository.On("GetTagBy", ctx, tags.UnknownSentinel().Key).Return(searchTagPtr(tags.UnknownSentinel()), nil)

	bus := expense.NewEventBus()
	repo := newDetectionRepository(mockSearchTagRepository, bus)

	result, err := repo.fromDynamo(ctx, dynamoItemFor(tags.UnknownSentinel().Key))

	assert.Equal(t, nil, err)
	assert.Equal(t, []tags.SearchTag{tags.UnknownSentinel()}, result.Tags)
	assertNoReclassifyEventFired(t, bus)
}

// A record that keeps a live tag but loses another resolves to [live, UNKNOWN]:
// the deleted reference must still be durably reclassified.
func TestPartialTagDeletionIsReclassified(t *testing.T) {
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockSearchTagRepository.On("GetTagBy", ctx, "super-market").Return(&tags.SearchTag{Key: "super-market", Value: "super-market"}, nil)
	mockSearchTagRepository.On("GetTagBy", ctx, "deleted-tag").Return(searchTagPtr(tags.UnknownSentinel()), nil)

	bus := expense.NewEventBus()
	repo := newDetectionRepository(mockSearchTagRepository, bus)

	result, err := repo.fromDynamo(ctx, dynamoItemFor("super-market,deleted-tag"))

	assert.Equal(t, nil, err)
	assert.Equal(t, []tags.SearchTag{{Key: "super-market", Value: "super-market"}, tags.UnknownSentinel()}, result.Tags)
	event := awaitReclassifyEvent(t, bus)
	assert.Equal(t, []tags.SearchTag{{Key: "super-market", Value: "super-market"}, tags.UnknownSentinel()}, event.Payload.Tags)
}

// Duplicate stored keys that both resolve to a live tag are not a deletion, so
// no reclassify event is fired (the read output keeps the duplicates).
func TestDuplicateStoredLiveTagsAreNotTreatedAsDeletion(t *testing.T) {
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockSearchTagRepository.On("GetTagBy", ctx, "dinner").Return(&tags.SearchTag{Key: "dinner", Value: "dinner"}, nil)

	bus := expense.NewEventBus()
	repo := newDetectionRepository(mockSearchTagRepository, bus)

	result, err := repo.fromDynamo(ctx, dynamoItemFor("dinner,dinner"))

	assert.Equal(t, nil, err)
	assert.Equal(t, []tags.SearchTag{{Key: "dinner", Value: "dinner"}, {Key: "dinner", Value: "dinner"}}, result.Tags)
	assertNoReclassifyEventFired(t, bus)
}

func newDetectionRepository(searchTagRepository tags.SearchTagRepository, bus *expense.InternalEventBus) *DynamoDbBudgetExpenseRepository {
	return NewDynamoDbBudgetExpenseRepository(TableName, client, new(DynamoDbBudgetExpenseIdProviderMock), searchTagRepository, bus).(*DynamoDbBudgetExpenseRepository)
}

func searchTagPtr(tag tags.SearchTag) *tags.SearchTag {
	return &tag
}

// dynamoItemFor builds a raw DynamoDB item with the given comma-separated tag
// keys so fromDynamo can be exercised without persisting through LocalStack.
func dynamoItemFor(tag string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"budget_id":        &types.AttributeValueMemberS{Value: "A_BUDGET_ID"},
		"user_name":        &types.AttributeValueMemberS{Value: "testuser"},
		"transaction_date": &types.AttributeValueMemberS{Value: "2018-01-01"},
		"amount":           &types.AttributeValueMemberS{Value: "10.50"},
		"note":             &types.AttributeValueMemberS{Value: "NOTE"},
		"tag":              &types.AttributeValueMemberS{Value: tag},
	}
}

func awaitReclassifyEvent(t *testing.T, bus *expense.InternalEventBus) expense.InternalEvent {
	t.Helper()
	select {
	case event := <-bus.Events():
		return event
	case <-time.After(2 * time.Second):
		t.Fatalf("expected a durable-reclassify event but none was fired")
		return expense.InternalEvent{}
	}
}

func assertNoReclassifyEventFired(t *testing.T, bus *expense.InternalEventBus) {
	t.Helper()
	select {
	case <-bus.Events():
		t.Fatalf("expected no durable-reclassify event but one was fired")
	case <-time.After(200 * time.Millisecond):
	}
}
