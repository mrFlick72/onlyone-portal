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

// The fromDynamo tests exercise the decoder directly with a hand-built item, no
// LocalStack round-trip, and assert the deleted-tag flag it returns (fromDynamo
// itself has no side effect — the caller decides whether to reclassify). A
// stored key that resolves to the UNKNOWN sentinel is a deletion only if the
// stored key itself was not UNKNOWN (a record stored with the sentinel as its
// key is the create-time default, not a deletion).

func TestFromDynamoFlagsASingleDeletedTag(t *testing.T) {
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockSearchTagRepository.On("GetTagBy", ctx, "deleted-tag").Return(searchTagPtr(tags.UnknownSentinel()), nil)
	repo := newDetectionRepository(mockSearchTagRepository)

	result, needsReclassification, err := repo.fromDynamo(ctx, dynamoItemFor("deleted-tag"))

	assert.Equal(t, nil, err)
	assert.Equal(t, true, needsReclassification)
	assert.Equal(t, []tags.SearchTag{tags.UnknownSentinel()}, result.Tags)
}

func TestFromDynamoFlagsEveryTagRemovedAndKeepsReadOutputUndeduplicated(t *testing.T) {
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockSearchTagRepository.On("GetTagBy", ctx, "gone-a").Return(searchTagPtr(tags.UnknownSentinel()), nil)
	mockSearchTagRepository.On("GetTagBy", ctx, "gone-b").Return(searchTagPtr(tags.UnknownSentinel()), nil)
	repo := newDetectionRepository(mockSearchTagRepository)

	result, needsReclassification, err := repo.fromDynamo(ctx, dynamoItemFor("gone-a,gone-b"))

	assert.Equal(t, nil, err)
	assert.Equal(t, true, needsReclassification)
	assert.Equal(t, []tags.SearchTag{tags.UnknownSentinel(), tags.UnknownSentinel()}, result.Tags)
}

func TestFromDynamoDoesNotFlagAHealthyRecord(t *testing.T) {
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockSearchTagRepository.On("GetTagBy", ctx, "super-market").Return(&tags.SearchTag{Key: "super-market", Value: "super-market"}, nil)
	repo := newDetectionRepository(mockSearchTagRepository)

	result, needsReclassification, err := repo.fromDynamo(ctx, dynamoItemFor("super-market"))

	assert.Equal(t, nil, err)
	assert.Equal(t, false, needsReclassification)
	assert.Equal(t, []tags.SearchTag{{Key: "super-market", Value: "super-market"}}, result.Tags)
}

// A record stored with the UNKNOWN sentinel as its tag key (the create-time
// default) resolves to UNKNOWN but was not deleted, so it must not be flagged —
// otherwise every untagged record would re-Save and re-emit on every read.
func TestFromDynamoDoesNotFlagARecordStoredWithTheUnknownSentinel(t *testing.T) {
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockSearchTagRepository.On("GetTagBy", ctx, tags.UnknownSentinel().Key).Return(searchTagPtr(tags.UnknownSentinel()), nil)
	repo := newDetectionRepository(mockSearchTagRepository)

	result, needsReclassification, err := repo.fromDynamo(ctx, dynamoItemFor(tags.UnknownSentinel().Key))

	assert.Equal(t, nil, err)
	assert.Equal(t, false, needsReclassification)
	assert.Equal(t, []tags.SearchTag{tags.UnknownSentinel()}, result.Tags)
}

// A record that keeps a live tag but loses another resolves to [live, UNKNOWN]:
// the deleted reference must still be flagged for reclassification.
func TestFromDynamoFlagsPartialTagDeletion(t *testing.T) {
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockSearchTagRepository.On("GetTagBy", ctx, "super-market").Return(&tags.SearchTag{Key: "super-market", Value: "super-market"}, nil)
	mockSearchTagRepository.On("GetTagBy", ctx, "deleted-tag").Return(searchTagPtr(tags.UnknownSentinel()), nil)
	repo := newDetectionRepository(mockSearchTagRepository)

	result, needsReclassification, err := repo.fromDynamo(ctx, dynamoItemFor("super-market,deleted-tag"))

	assert.Equal(t, nil, err)
	assert.Equal(t, true, needsReclassification)
	assert.Equal(t, []tags.SearchTag{{Key: "super-market", Value: "super-market"}, tags.UnknownSentinel()}, result.Tags)
}

// Duplicate stored keys that both resolve to a live tag are not a deletion.
func TestFromDynamoDoesNotFlagDuplicateLiveTags(t *testing.T) {
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockSearchTagRepository.On("GetTagBy", ctx, "dinner").Return(&tags.SearchTag{Key: "dinner", Value: "dinner"}, nil)
	repo := newDetectionRepository(mockSearchTagRepository)

	result, needsReclassification, err := repo.fromDynamo(ctx, dynamoItemFor("dinner,dinner"))

	assert.Equal(t, nil, err)
	assert.Equal(t, false, needsReclassification)
	assert.Equal(t, []tags.SearchTag{{Key: "dinner", Value: "dinner"}, {Key: "dinner", Value: "dinner"}}, result.Tags)
}

// The reclassify payload collapses repeated keys so a fully-deleted record
// persists a single UNKNOWN rather than repeated sentinels.
func TestDistinctTagsCollapsesRepeatedKeys(t *testing.T) {
	superMarket := tags.SearchTag{Key: "super-market", Value: "super-market"}

	got := distinctTags([]tags.SearchTag{tags.UnknownSentinel(), tags.UnknownSentinel(), superMarket})

	assert.Equal(t, []tags.SearchTag{tags.UnknownSentinel(), superMarket}, got)
}

// Regression for review finding #1: FindFor is the ownership pre-check used by
// update and delete, never a standalone read. Reading a deleted-tag record via
// FindFor must resolve UNKNOWN in place but publish NOTHING — otherwise the async
// re-Save of the pre-mutation snapshot would resurrect a just-deleted expense or
// clobber an in-flight update.
func TestFindForDoesNotPublishReclassification(t *testing.T) {
	guardCtx := testutils.NewStubbedContextWith("delete-guard-user")
	idProvider := &DynamoDbBudgetExpenseIdProvider{SaltGenerator: func() string { return uuid.New().String() }}
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockSearchTagRepository.On("GetTagBy", guardCtx, "deleted-tag").Return(searchTagPtr(tags.UnknownSentinel()), nil)

	bus := expense.NewEventBus()
	repo := NewDynamoDbBudgetExpenseRepository(TableName, client, idProvider, mockSearchTagRepository, bus).(*DynamoDbBudgetExpenseRepository)

	stored := expense.BudgetExpense{
		UserName: "delete-guard-user",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tags:     []tags.SearchTag{{Key: "deleted-tag", Value: "deleted-tag"}},
	}
	if err := repo.Save(guardCtx, &stored); err != nil {
		t.Fatalf("failed to save fixture: %v", err)
	}

	found, err := repo.FindFor(guardCtx, stored.Id)

	assert.Equal(t, nil, err)
	assert.Equal(t, []tags.SearchTag{tags.UnknownSentinel()}, found.Tags)
	assertNoReclassifyEventFired(t, bus)
}

// The genuine read path still reclassifies: a deleted-tag reference found during
// a range read publishes a reclassification event with the deduplicated payload.
func TestFindByDateRangePublishesReclassificationForADeletedTag(t *testing.T) {
	readCtx := testutils.NewStubbedContextWith("reclassify-read-user")
	idProvider := &DynamoDbBudgetExpenseIdProvider{SaltGenerator: func() string { return uuid.New().String() }}
	mockSearchTagRepository := new(tags.SearchTagRepositoryMock)
	mockSearchTagRepository.On("GetTagBy", readCtx, "deleted-tag").Return(searchTagPtr(tags.UnknownSentinel()), nil)

	bus := expense.NewEventBus()
	repo := NewDynamoDbBudgetExpenseRepository(TableName, client, idProvider, mockSearchTagRepository, bus).(*DynamoDbBudgetExpenseRepository)

	stored := expense.BudgetExpense{
		UserName: "reclassify-read-user",
		Date:     testutils.SafeDateFor("05/03/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
		Tags:     []tags.SearchTag{{Key: "deleted-tag", Value: "deleted-tag"}},
	}
	if err := repo.Save(readCtx, &stored); err != nil {
		t.Fatalf("failed to save fixture: %v", err)
	}

	result, err := repo.FindByDateRange(readCtx, testutils.SafeDateFor("01/03/2024"), testutils.SafeDateFor("31/03/2024"), []tags.SearchTagKey{})

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(result))
	event := awaitReclassifyEvent(t, bus)
	assert.Equal(t, stored.Id, event.Payload.Id)
	assert.Equal(t, []tags.SearchTag{tags.UnknownSentinel()}, event.Payload.Tags)
}

func newDetectionRepository(searchTagRepository tags.SearchTagRepository) *DynamoDbBudgetExpenseRepository {
	return NewDynamoDbBudgetExpenseRepository(TableName, client, new(DynamoDbBudgetExpenseIdProviderMock), searchTagRepository, expense.NewEventBus()).(*DynamoDbBudgetExpenseRepository)
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
