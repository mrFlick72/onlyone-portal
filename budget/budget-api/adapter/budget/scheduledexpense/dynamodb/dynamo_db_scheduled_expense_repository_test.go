//go:build test

package dynamodb

import (
	"os"
	"testing"

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
