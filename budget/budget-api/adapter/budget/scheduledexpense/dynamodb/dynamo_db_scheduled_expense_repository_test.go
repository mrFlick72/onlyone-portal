//go:build test

package dynamodb

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/budget/scheduledexpense"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestMain(m *testing.M) {
	setupTestDynamoDBTable()
	code := m.Run()
	teardownTestDynamoDBTable()
	os.Exit(code)
}

func TestSaveAndFindAllReturnsAnUntaggedScheduledExpense(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	userCtx := testutils.NewStubbedContextWith("user-untagged")

	input := scheduledexpense.ScheduledExpense{
		UserName:    "user-untagged",
		Description: "Rent",
		Amount:      testutils.SafeMoneyFor("1200.00"),
		Notes:       "Monthly rent",
		Day:         5,
		Status:      scheduledexpense.StatusActive,
	}
	idProviderMock.On("GenerateIdFor", &input).Return("RENT_ID")

	if err := repo.Save(userCtx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	result, err := repo.FindAll(userCtx)
	if err != nil {
		t.Fatalf("Error finding scheduled expenses: %v", err)
	}

	expected := scheduledexpense.ScheduledExpense{
		Id:          "RENT_ID",
		UserName:    "user-untagged",
		Description: "Rent",
		Amount:      testutils.SafeMoneyFor("1200.00"),
		Notes:       "Monthly rent",
		Tags:        []tags.SearchTag{tags.UnknownSentinel()},
		Day:         5,
		Status:      scheduledexpense.StatusActive,
	}

	assert.Equal(t, 1, len(result))
	assert.Equal(t, expected, result[0])
}

func TestSaveAndFindAllResolvesTaggedScheduledExpense(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	tagRepositoryMock := new(tags.SearchTagRepositoryMock)
	repo := newScheduledExpenseRepositoryWith(idProviderMock, tagRepositoryMock)
	userCtx := testutils.NewStubbedContextWith("user-tagged")

	housing := tags.SearchTag{Key: "housing", Value: "Housing"}
	tagRepositoryMock.On("GetTagBy", userCtx, "housing").Return(&housing, nil)

	input := scheduledexpense.ScheduledExpense{
		UserName:    "user-tagged",
		Description: "Rent tagged",
		Amount:      testutils.SafeMoneyFor("1200.00"),
		Notes:       "Monthly rent",
		Tags:        []tags.SearchTag{housing},
		Day:         5,
		Status:      scheduledexpense.StatusActive,
	}
	idProviderMock.On("GenerateIdFor", &input).Return("RENT_TAGGED_ID")

	if err := repo.Save(userCtx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	result, err := repo.FindAll(userCtx)
	if err != nil {
		t.Fatalf("Error finding scheduled expense: %v", err)
	}

	assert.Equal(t, 1, len(result))
	assert.Equal(t, []tags.SearchTag{housing}, result[0].Tags)
	tagRepositoryMock.AssertCalled(t, "GetTagBy", userCtx, "housing")
}

func TestSaveAndFindAllRoundTripsOptionalFieldsWhenSet(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	userCtx := testutils.NewStubbedContextWith("user-optional-set")

	month := 3
	endDate := testutils.SafeDateFor("31/12/2026")
	lastEvaluated := testutils.SafeDateFor("15/03/2026")

	input := scheduledexpense.ScheduledExpense{
		UserName:          "user-optional-set",
		Description:       "Insurance",
		Amount:            testutils.SafeMoneyFor("450.00"),
		Notes:             "Yearly insurance",
		Day:               15,
		Month:             &month,
		EndDate:           &endDate,
		Status:            scheduledexpense.StatusActive,
		LastEvaluatedDate: &lastEvaluated,
	}
	idProviderMock.On("GenerateIdFor", &input).Return("INSURANCE_ID")

	if err := repo.Save(userCtx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	result, err := repo.FindAll(userCtx)
	if err != nil {
		t.Fatalf("Error finding scheduled expenses: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected 1 scheduled expense, got %d", len(result))
	}
	found := result[0]

	if found.Month == nil || *found.Month != month {
		t.Fatalf("Expected Month %d, got %v", month, found.Month)
	}
	if found.EndDate == nil || found.EndDate.GetIsoFormattedDate() != endDate.GetIsoFormattedDate() {
		t.Fatalf("Expected EndDate %v, got %v", endDate, found.EndDate)
	}
	if found.LastEvaluatedDate == nil || found.LastEvaluatedDate.GetIsoFormattedDate() != lastEvaluated.GetIsoFormattedDate() {
		t.Fatalf("Expected LastEvaluatedDate %v, got %v", lastEvaluated, found.LastEvaluatedDate)
	}
}

func TestSaveWithoutOptionalFieldsLeavesThemNil(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	userCtx := testutils.NewStubbedContextWith("user-optional-unset")

	input := scheduledexpense.ScheduledExpense{
		UserName:    "user-optional-unset",
		Description: "Subscription",
		Amount:      testutils.SafeMoneyFor("9.99"),
		Day:         1,
		Status:      scheduledexpense.StatusActive,
	}
	idProviderMock.On("GenerateIdFor", &input).Return("SUBSCRIPTION_ID")

	if err := repo.Save(userCtx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	result, err := repo.FindAll(userCtx)
	if err != nil {
		t.Fatalf("Error finding scheduled expenses: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected 1 scheduled expense, got %d", len(result))
	}
	found := result[0]

	assert.Equal(t, true, found.Month == nil)
	assert.Equal(t, true, found.EndDate == nil)
	assert.Equal(t, true, found.LastEvaluatedDate == nil)
}

// TestSaveDerivesUserNameFromContextNotFromTheStruct guards the ownership
// boundary Save enforces: the partition key must come from the authenticated
// caller (ctx), never from a caller-supplied struct field, the same
// enforcement revenue's Save has (see budget-api/CLAUDE.md "Ownership
// enforced in two layers — keep both"). Without this, a struct built from an
// unpopulated representation (UserName == "") would write into a phantom
// partition instead of failing or scoping correctly.
func TestSaveDerivesUserNameFromContextNotFromTheStruct(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	userCtx := testutils.NewStubbedContextWith("user-ctx-wins")

	input := scheduledexpense.ScheduledExpense{
		UserName:    "SOMEONE_ELSE",
		Description: "Rent",
		Amount:      testutils.SafeMoneyFor("1.00"),
		Day:         1,
		Status:      scheduledexpense.StatusActive,
	}
	idProviderMock.On("GenerateIdFor", &input).Return("CTX_ID")

	if err := repo.Save(userCtx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	assert.Equal(t, "user-ctx-wins", input.UserName)

	result, err := repo.FindAll(userCtx)
	if err != nil {
		t.Fatalf("Error finding scheduled expenses: %v", err)
	}
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "user-ctx-wins", result[0].UserName)
}

func TestFindForReturnsTheScheduledExpenseWhenItExists(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	userCtx := testutils.NewStubbedContextWith("user-find-for")

	input := scheduledexpense.ScheduledExpense{
		UserName:    "user-find-for",
		Description: "Rent",
		Amount:      testutils.SafeMoneyFor("1200.00"),
		Day:         5,
		Status:      scheduledexpense.StatusActive,
	}
	idProviderMock.On("GenerateIdFor", &input).Return("FIND_FOR_ID")
	if err := repo.Save(userCtx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	found, err := repo.FindFor(userCtx, "FIND_FOR_ID")

	if err != nil {
		t.Fatalf("Error finding scheduled expense: %v", err)
	}
	if found == nil {
		t.Fatalf("Expected a scheduled expense, got nil")
	}
	assert.Equal(t, "FIND_FOR_ID", found.Id)
	assert.Equal(t, "Rent", found.Description)
}

func TestFindForReturnsNilWithoutErrorWhenNotFound(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	userCtx := testutils.NewStubbedContextWith("user-find-for-missing")

	found, err := repo.FindFor(userCtx, "DOES_NOT_EXIST")

	assert.Equal(t, nil, err)
	if found != nil {
		t.Fatalf("Expected nil, got %+v", found)
	}
}

func TestFindForDoesNotLeakAnotherUsersScheduledExpense(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	ownerCtx := testutils.NewStubbedContextWith("user-find-for-owner")
	otherCtx := testutils.NewStubbedContextWith("user-find-for-other")

	input := scheduledexpense.ScheduledExpense{
		UserName:    "user-find-for-owner",
		Description: "Rent",
		Amount:      testutils.SafeMoneyFor("1200.00"),
		Day:         5,
		Status:      scheduledexpense.StatusActive,
	}
	idProviderMock.On("GenerateIdFor", &input).Return("OWNER_ONLY_ID")
	if err := repo.Save(ownerCtx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	found, err := repo.FindFor(otherCtx, "OWNER_ONLY_ID")

	assert.Equal(t, nil, err)
	if found != nil {
		t.Fatalf("Expected FindFor under another user's context to find nothing, got %+v", found)
	}
}

func TestFindAllOnlyReturnsTheCurrentUsersScheduledExpenses(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	mineCtx := testutils.NewStubbedContextWith("user-scoped-mine")
	theirsCtx := testutils.NewStubbedContextWith("user-scoped-theirs")

	mine := scheduledexpense.ScheduledExpense{UserName: "user-scoped-mine", Description: "Mine", Amount: testutils.SafeMoneyFor("1.00"), Day: 1, Status: scheduledexpense.StatusActive}
	idProviderMock.On("GenerateIdFor", &mine).Return("MINE_ID")
	if err := repo.Save(mineCtx, &mine); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	theirs := scheduledexpense.ScheduledExpense{UserName: "user-scoped-theirs", Description: "Theirs", Amount: testutils.SafeMoneyFor("1.00"), Day: 1, Status: scheduledexpense.StatusActive}
	idProviderMock.On("GenerateIdFor", &theirs).Return("THEIRS_ID")
	if err := repo.Save(theirsCtx, &theirs); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	result, err := repo.FindAll(mineCtx)
	if err != nil {
		t.Fatalf("Error finding scheduled expenses: %v", err)
	}

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "MINE_ID", result[0].Id)
}

func TestDeleteRemovesTheScheduledExpense(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	userCtx := testutils.NewStubbedContextWith("user-delete")

	input := scheduledexpense.ScheduledExpense{
		UserName:    "user-delete",
		Description: "Rent",
		Amount:      testutils.SafeMoneyFor("1200.00"),
		Day:         5,
		Status:      scheduledexpense.StatusActive,
	}
	idProviderMock.On("GenerateIdFor", &input).Return("DELETE_ID")
	if err := repo.Save(userCtx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	err := repo.Delete(userCtx, "DELETE_ID")

	assert.Equal(t, nil, err)
	found, err := repo.FindFor(userCtx, "DELETE_ID")
	assert.Equal(t, nil, err)
	if found != nil {
		t.Fatalf("Expected the scheduled expense to be gone, got %+v", found)
	}
}

// The attribute_exists(id) condition turns a delete of a row that isn't
// there (e.g. removed between the action's FindFor and this call) into a
// not-found, instead of DynamoDB's silent no-op success.
func TestDeleteReturnsNotFoundWhenTheScheduledExpenseDoesNotExist(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	userCtx := testutils.NewStubbedContextWith("user-delete-missing")

	err := repo.Delete(userCtx, "DOES_NOT_EXIST")

	assert.Equal(t, scheduledexpense.ErrScheduledExpenseNotFound, err)
}

func TestDeleteUnderAnotherUsersContextLeavesTheOwnersScheduledExpense(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	ownerCtx := testutils.NewStubbedContextWith("user-delete-owner")
	otherCtx := testutils.NewStubbedContextWith("user-delete-other")

	input := scheduledexpense.ScheduledExpense{
		UserName:    "user-delete-owner",
		Description: "Rent",
		Amount:      testutils.SafeMoneyFor("1200.00"),
		Day:         5,
		Status:      scheduledexpense.StatusActive,
	}
	idProviderMock.On("GenerateIdFor", &input).Return("OWNER_DELETE_ID")
	if err := repo.Save(ownerCtx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	err := repo.Delete(otherCtx, "OWNER_DELETE_ID")

	assert.Equal(t, scheduledexpense.ErrScheduledExpenseNotFound, err)
	found, err := repo.FindFor(ownerCtx, "OWNER_DELETE_ID")
	assert.Equal(t, nil, err)
	if found == nil {
		t.Fatalf("Expected the owner's scheduled expense to survive another user's delete")
	}
}

// UpdateStatus touches only status + last_evaluated_date: description, amount
// and the stored tag keys must survive untouched (a full-item Save of a
// FindFor-resolved record could rewrite tag keys, e.g. to UNKNOWN).
func TestUpdateStatusSetsOnlyStatusAndLastEvaluatedDate(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	tagRepositoryMock := new(tags.SearchTagRepositoryMock)
	repo := newScheduledExpenseRepositoryWith(idProviderMock, tagRepositoryMock)
	userCtx := testutils.NewStubbedContextWith("user-update-status")

	housing := tags.SearchTag{Key: "housing", Value: "Housing"}
	tagRepositoryMock.On("GetTagBy", userCtx, "housing").Return(&housing, nil)

	input := scheduledexpense.ScheduledExpense{
		UserName:    "user-update-status",
		Description: "Rent",
		Amount:      testutils.SafeMoneyFor("1200.00"),
		Notes:       "Monthly rent",
		Tags:        []tags.SearchTag{housing},
		Day:         5,
		Status:      scheduledexpense.StatusActive,
	}
	idProviderMock.On("GenerateIdFor", &input).Return("UPDATE_STATUS_ID")
	if err := repo.Save(userCtx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	today := testutils.SafeDateFor("24/09/2026")
	err := repo.UpdateStatus(userCtx, "UPDATE_STATUS_ID", scheduledexpense.StatusPaused, today)

	assert.Equal(t, nil, err)
	found, err := repo.FindFor(userCtx, "UPDATE_STATUS_ID")
	if err != nil || found == nil {
		t.Fatalf("Expected the scheduled expense, got %+v / %v", found, err)
	}
	assert.Equal(t, scheduledexpense.StatusPaused, found.Status)
	if found.LastEvaluatedDate == nil {
		t.Fatalf("Expected LastEvaluatedDate to be stamped")
	}
	assert.Equal(t, "2026-09-24", found.LastEvaluatedDate.GetIsoFormattedDate())
	assert.Equal(t, "Rent", found.Description)
	assert.Equal(t, "1200.00", found.Amount.StringifyAmount())
	assert.Equal(t, "Monthly rent", found.Notes)
	assert.Equal(t, []tags.SearchTag{housing}, found.Tags)
}

func TestUpdateStatusReturnsNotFoundWhenTheScheduledExpenseDoesNotExist(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	userCtx := testutils.NewStubbedContextWith("user-update-status-missing")

	err := repo.UpdateStatus(userCtx, "DOES_NOT_EXIST", scheduledexpense.StatusPaused, testutils.SafeDateFor("24/09/2026"))

	assert.Equal(t, scheduledexpense.ErrScheduledExpenseNotFound, err)
	// ...and the condition stops UpdateItem's upsert from creating a stub row.
	found, err := repo.FindFor(userCtx, "DOES_NOT_EXIST")
	assert.Equal(t, nil, err)
	if found != nil {
		t.Fatalf("Expected no row to be created, got %+v", found)
	}
}

func TestUpdateStatusUnderAnotherUsersContextLeavesTheOwnersScheduledExpense(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	ownerCtx := testutils.NewStubbedContextWith("user-update-status-owner")
	otherCtx := testutils.NewStubbedContextWith("user-update-status-other")

	input := scheduledexpense.ScheduledExpense{
		UserName:    "user-update-status-owner",
		Description: "Rent",
		Amount:      testutils.SafeMoneyFor("1200.00"),
		Day:         5,
		Status:      scheduledexpense.StatusActive,
	}
	idProviderMock.On("GenerateIdFor", &input).Return("OWNER_STATUS_ID")
	if err := repo.Save(ownerCtx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	err := repo.UpdateStatus(otherCtx, "OWNER_STATUS_ID", scheduledexpense.StatusPaused, testutils.SafeDateFor("24/09/2026"))

	assert.Equal(t, scheduledexpense.ErrScheduledExpenseNotFound, err)
	found, err := repo.FindFor(ownerCtx, "OWNER_STATUS_ID")
	if err != nil || found == nil {
		t.Fatalf("Expected the owner's scheduled expense, got %+v / %v", found, err)
	}
	assert.Equal(t, scheduledexpense.StatusActive, found.Status)
	if found.LastEvaluatedDate != nil {
		t.Fatalf("Expected the owner's LastEvaluatedDate untouched, got %v", found.LastEvaluatedDate)
	}
}

// findActiveById picks this test's rows out of a table-wide scan (the table is
// shared by every test in the package).
func findActiveById(t *testing.T, repo *DynamoDbScheduledExpenseRepository, ids ...string) map[string]scheduledexpense.ScheduledExpense {
	t.Helper()
	all, err := repo.FindAllActive(context.Background())
	if err != nil {
		t.Fatalf("Error scanning active scheduled expenses: %v", err)
	}
	wanted := map[string]bool{}
	for _, id := range ids {
		wanted[id] = true
	}
	found := map[string]scheduledexpense.ScheduledExpense{}
	for _, se := range all {
		if wanted[se.Id] {
			found[se.Id] = se
		}
	}
	return found
}

// FindAllActive spans every user's partition, returns only ACTIVE rows, and
// reads the tag names stored at save time — the tag repository mock has no
// expectations, so any tag-api resolution would fail the test.
func TestFindAllActiveReturnsEveryUsersActiveDefinitionsWithStoredTagNames(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	aliceCtx := testutils.NewStubbedContextWith("user-scan-alice")
	bobCtx := testutils.NewStubbedContextWith("user-scan-bob")

	alice := scheduledexpense.ScheduledExpense{
		Description: "Rent", Amount: testutils.SafeMoneyFor("1200.00"), Notes: "Monthly rent",
		Tags: []tags.SearchTag{{Key: "housing", Value: "Housing"}, {Key: "fixed", Value: "Fixed costs"}},
		Day:  5, Status: scheduledexpense.StatusActive,
	}
	bob := scheduledexpense.ScheduledExpense{
		Description: "Gym", Amount: testutils.SafeMoneyFor("40.00"),
		Tags: []tags.SearchTag{{Key: "sport", Value: "Sport"}},
		Day:  1, Status: scheduledexpense.StatusActive,
	}
	bobPaused := scheduledexpense.ScheduledExpense{
		Description: "Magazine", Amount: testutils.SafeMoneyFor("5.00"),
		Tags: []tags.SearchTag{{Key: "leisure", Value: "Leisure"}},
		Day:  1, Status: scheduledexpense.StatusPaused,
	}
	idProviderMock.On("GenerateIdFor", &alice).Return("SCAN_ALICE_ID")
	idProviderMock.On("GenerateIdFor", &bob).Return("SCAN_BOB_ID")
	idProviderMock.On("GenerateIdFor", &bobPaused).Return("SCAN_BOB_PAUSED_ID")
	for _, save := range []struct {
		ctx context.Context
		se  *scheduledexpense.ScheduledExpense
	}{{aliceCtx, &alice}, {bobCtx, &bob}, {bobCtx, &bobPaused}} {
		if err := repo.Save(save.ctx, save.se); err != nil {
			t.Fatalf("Expected nil error on save, got %v", err)
		}
	}

	found := findActiveById(t, repo, "SCAN_ALICE_ID", "SCAN_BOB_ID", "SCAN_BOB_PAUSED_ID")

	assert.Equal(t, 2, len(found))
	assert.Equal(t, "user-scan-alice", found["SCAN_ALICE_ID"].UserName)
	assert.Equal(t, "Rent", found["SCAN_ALICE_ID"].Description)
	assert.Equal(t, "Monthly rent", found["SCAN_ALICE_ID"].Notes)
	assert.Equal(t, []tags.SearchTag{{Key: "housing", Value: "Housing"}, {Key: "fixed", Value: "Fixed costs"}}, found["SCAN_ALICE_ID"].Tags)
	assert.Equal(t, "user-scan-bob", found["SCAN_BOB_ID"].UserName)
	assert.Equal(t, []tags.SearchTag{{Key: "sport", Value: "Sport"}}, found["SCAN_BOB_ID"].Tags)
}

// Rows saved before tag names were stored still scan, with empty names.
func TestFindAllActiveToleratesRowsWithoutStoredTagNames(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)

	_, err := client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item: map[string]types.AttributeValue{
			"user_name":   &types.AttributeValueMemberS{Value: "user-scan-legacy"},
			"id":          &types.AttributeValueMemberS{Value: "SCAN_LEGACY_ID"},
			"description": &types.AttributeValueMemberS{Value: "Legacy"},
			"amount":      &types.AttributeValueMemberS{Value: "10.00"},
			"notes":       &types.AttributeValueMemberS{Value: ""},
			"tag":         &types.AttributeValueMemberS{Value: "housing"},
			"day":         &types.AttributeValueMemberN{Value: "3"},
			"status":      &types.AttributeValueMemberS{Value: "ACTIVE"},
		},
	})
	if err != nil {
		t.Fatalf("Error seeding legacy row: %v", err)
	}

	found := findActiveById(t, repo, "SCAN_LEGACY_ID")

	assert.Equal(t, []tags.SearchTag{{Key: "housing", Value: ""}}, found["SCAN_LEGACY_ID"].Tags)
}

// The scan follows LastEvaluatedKey across pages.
func TestFindAllActiveFollowsPagination(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	repo.scanPageLimit = 1
	userCtx := testutils.NewStubbedContextWith("user-scan-pages")

	ids := []string{"SCAN_PAGE_1", "SCAN_PAGE_2", "SCAN_PAGE_3"}
	for i, id := range ids {
		se := scheduledexpense.ScheduledExpense{
			Description: fmt.Sprintf("Page %d", i), Amount: testutils.SafeMoneyFor("1.00"),
			Day: 1, Status: scheduledexpense.StatusActive,
		}
		idProviderMock.On("GenerateIdFor", &se).Return(id).Once()
		if err := repo.Save(userCtx, &se); err != nil {
			t.Fatalf("Expected nil error on save, got %v", err)
		}
	}

	found := findActiveById(t, repo, ids...)

	assert.Equal(t, 3, len(found))
}

func TestAdvanceLastEvaluatedDateSetsOnlyThatField(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	userCtx := testutils.NewStubbedContextWith("user-advance")

	input := scheduledexpense.ScheduledExpense{
		Description: "Rent", Amount: testutils.SafeMoneyFor("1200.00"),
		Tags: []tags.SearchTag{{Key: "housing", Value: "Housing"}},
		Day:  5, Status: scheduledexpense.StatusActive,
	}
	idProviderMock.On("GenerateIdFor", &input).Return("ADVANCE_ID")
	if err := repo.Save(userCtx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	err := repo.AdvanceLastEvaluatedDate(userCtx, "ADVANCE_ID", testutils.SafeDateFor("24/09/2026"))

	assert.Equal(t, nil, err)
	found := findActiveById(t, repo, "ADVANCE_ID")["ADVANCE_ID"]
	if found.LastEvaluatedDate == nil {
		t.Fatalf("Expected LastEvaluatedDate to be set")
	}
	assert.Equal(t, "2026-09-24", found.LastEvaluatedDate.GetIsoFormattedDate())
	assert.Equal(t, scheduledexpense.StatusActive, found.Status)
	assert.Equal(t, "Rent", found.Description)
	assert.Equal(t, []tags.SearchTag{{Key: "housing", Value: "Housing"}}, found.Tags)
}

// A definition deleted mid-run: not found, and no stub row upserted (a stub
// would lack description/amount and break every read of that user's list).
func TestAdvanceLastEvaluatedDateOnADeletedDefinitionIsNotFoundAndCreatesNothing(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	userCtx := testutils.NewStubbedContextWith("user-advance-deleted")

	err := repo.AdvanceLastEvaluatedDate(userCtx, "ADVANCE_DELETED_ID", testutils.SafeDateFor("24/09/2026"))

	assert.Equal(t, scheduledexpense.ErrScheduledExpenseNotFound, err)
	found, err := repo.FindFor(userCtx, "ADVANCE_DELETED_ID")
	assert.Equal(t, nil, err)
	if found != nil {
		t.Fatalf("Expected no stub row, got %+v", found)
	}
}

// A definition paused mid-run stops advancing: pause owns LastEvaluatedDate
// from then on (ADR 0005).
func TestAdvanceLastEvaluatedDateOnAPausedDefinitionIsNotFoundAndLeavesItUntouched(t *testing.T) {
	idProviderMock := new(DynamoDbScheduledExpenseIdProviderMock)
	repo := newScheduledExpenseRepository(idProviderMock)
	userCtx := testutils.NewStubbedContextWith("user-advance-paused")

	input := scheduledexpense.ScheduledExpense{
		Description: "Rent", Amount: testutils.SafeMoneyFor("1200.00"),
		Day: 5, Status: scheduledexpense.StatusActive,
	}
	idProviderMock.On("GenerateIdFor", &input).Return("ADVANCE_PAUSED_ID")
	if err := repo.Save(userCtx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}
	if err := repo.UpdateStatus(userCtx, "ADVANCE_PAUSED_ID", scheduledexpense.StatusPaused, testutils.SafeDateFor("20/09/2026")); err != nil {
		t.Fatalf("Expected nil error on pause, got %v", err)
	}

	err := repo.AdvanceLastEvaluatedDate(userCtx, "ADVANCE_PAUSED_ID", testutils.SafeDateFor("24/09/2026"))

	assert.Equal(t, scheduledexpense.ErrScheduledExpenseNotFound, err)
	tagRepositoryMock := new(tags.SearchTagRepositoryMock)
	tagRepositoryMock.On("GetTagBy", userCtx, "UNKNOWN").Return(&tags.SearchTag{Key: "UNKNOWN", Value: "UNKNOWN"}, nil).Maybe()
	found, err := newScheduledExpenseRepositoryWith(idProviderMock, tagRepositoryMock).FindFor(userCtx, "ADVANCE_PAUSED_ID")
	if err != nil || found == nil {
		t.Fatalf("Expected the paused definition, got %+v / %v", found, err)
	}
	assert.Equal(t, "2026-09-20", found.LastEvaluatedDate.GetIsoFormattedDate())
}
