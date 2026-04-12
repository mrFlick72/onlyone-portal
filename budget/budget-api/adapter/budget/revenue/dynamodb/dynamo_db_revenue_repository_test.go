package dynamodb

import (
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/budget/revenue"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestMain(m *testing.M) {
	setupTestDynamoDBTable()
	code := m.Run()
	teardownTestDynamoDBTable()
	os.Exit(code)
}

func TestSaveANewRevenue(t *testing.T) {
	idProviderMock := new(DynamoDbRevenueIdProviderMock)
	repo := newRevenueRepository(idProviderMock)

	input := revenue.Revenue{
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
	}

	expected := revenue.Revenue{
		Id:       revenue.RevenueId("MjAyNF90ZXN0dXNlcg==-MV8xX0FfU0FMVA=="),
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
	}

	idProviderMock.On("GenerateIdFor", &input).Return(string(expected.Id))

	if err := repo.Save(ctx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	retrieved, err := repo.FindFor(ctx, expected.Id)
	if err != nil {
		t.Fatalf("Error retrieving revenue: %v", err)
	}

	idProviderMock.AssertCalled(t, "GenerateIdFor", &input)
	assert.Equal(t, expected, *retrieved)
}

func TestUpdateARevenue(t *testing.T) {
	idProviderMock := new(DynamoDbRevenueIdProviderMock)
	repo := newRevenueRepository(idProviderMock)

	input := revenue.Revenue{
		UserName: "testuser",
		Date:     testutils.SafeDateFor("02/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
	}
	idProviderMock.On("GenerateIdFor", &input).Return("MjAyNF90ZXN0dXNlcg==-MV8yX1NBTFQ=")

	if err := repo.Save(ctx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	updated := revenue.Revenue{
		Id:       revenue.RevenueId("MjAyNF90ZXN0dXNlcg==-MV8yX1NBTFQ="),
		UserName: "testuser",
		Date:     testutils.SafeDateFor("02/01/2024"),
		Amount:   testutils.SafeMoneyFor("99.99"),
		Note:     "UPDATED",
	}

	if err := repo.Save(ctx, &updated); err != nil {
		t.Fatalf("Expected nil error on update, got %v", err)
	}

	retrieved, err := repo.FindFor(ctx, updated.Id)
	if err != nil {
		t.Fatalf("Error retrieving revenue: %v", err)
	}

	assert.Equal(t, updated, *retrieved)
	idProviderMock.AssertNumberOfCalls(t, "GenerateIdFor", 1)
}

func TestUpdateARevenueFailsWhenNotOwner(t *testing.T) {
	idProviderMock := new(DynamoDbRevenueIdProviderMock)
	repo := newRevenueRepository(idProviderMock)

	input := revenue.Revenue{
		UserName: "testuser",
		Date:     testutils.SafeDateFor("03/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
	}
	idProviderMock.On("GenerateIdFor", &input).Return("MjAyNF90ZXN0dXNlcg==-MV8zX1NBTFQ=")

	if err := repo.Save(ctx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	attacker := revenue.Revenue{
		Id:       revenue.RevenueId("MjAyNF90ZXN0dXNlcg==-MV8zX1NBTFQ="),
		UserName: "anotheruser",
		Date:     testutils.SafeDateFor("03/01/2024"),
		Amount:   testutils.SafeMoneyFor("1.00"),
		Note:     "PWNED",
	}

	err := repo.Save(ctxAnotherUser, &attacker)
	assert.NotEqual(t, nil, err)
}

func TestFindNonExistentRevenueReturnsError(t *testing.T) {
	idProviderMock := new(DynamoDbRevenueIdProviderMock)
	repo := newRevenueRepository(idProviderMock)

	_, err := repo.FindFor(ctx, revenue.RevenueId("bm90LWZvdW5k-bm90LWZvdW5k"))
	assert.NotEqual(t, nil, err)
}

func TestFindRevenueOfOtherUserIsNotAllowed(t *testing.T) {
	idProviderMock := new(DynamoDbRevenueIdProviderMock)
	repo := newRevenueRepository(idProviderMock)

	input := revenue.Revenue{
		UserName: "testuser",
		Date:     testutils.SafeDateFor("04/01/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
	}
	id := revenue.RevenueId("MjAyNF90ZXN0dXNlcg==-MV80X1NBTFQ=")
	idProviderMock.On("GenerateIdFor", &input).Return(string(id))

	if err := repo.Save(ctx, &input); err != nil {
		t.Fatalf("Expected nil error on save, got %v", err)
	}

	retrieved, err := repo.FindFor(ctxAnotherUser, id)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, nil, retrieved)
}

func TestFindRevenueByDateRange(t *testing.T) {
	idProviderMock := new(DynamoDbRevenueIdProviderMock)
	repo := newRevenueRepository(idProviderMock)

	jan := revenue.Revenue{UserName: "testuser", Date: testutils.SafeDateFor("15/01/2025"), Amount: testutils.SafeMoneyFor("100.00"), Note: "JAN"}
	jun := revenue.Revenue{UserName: "testuser", Date: testutils.SafeDateFor("15/06/2025"), Amount: testutils.SafeMoneyFor("200.00"), Note: "JUN"}
	dec := revenue.Revenue{UserName: "testuser", Date: testutils.SafeDateFor("31/12/2025"), Amount: testutils.SafeMoneyFor("300.00"), Note: "DEC"}
	outOfRange := revenue.Revenue{UserName: "testuser", Date: testutils.SafeDateFor("15/06/2024"), Amount: testutils.SafeMoneyFor("999.00"), Note: "OLD"}

	idProviderMock.On("GenerateIdFor", &jan).Return("MjAyNV90ZXN0dXNlcg==-MV8xNV9TQTE=")
	idProviderMock.On("GenerateIdFor", &jun).Return("MjAyNV90ZXN0dXNlcg==-Nl8xNV9TQTI=")
	idProviderMock.On("GenerateIdFor", &dec).Return("MjAyNV90ZXN0dXNlcg==-MTJfMzFfU0Ez")
	idProviderMock.On("GenerateIdFor", &outOfRange).Return("MjAyNF90ZXN0dXNlcg==-Nl8xNV9TQTQ=")

	for _, r := range []*revenue.Revenue{&jan, &jun, &dec, &outOfRange} {
		if err := repo.Save(ctx, r); err != nil {
			t.Fatalf("save failed: %v", err)
		}
	}

	start := testutils.SafeDateFor("01/01/2025")
	end := testutils.SafeDateFor("31/12/2025")
	result, err := repo.FindByDateRange(ctx, start, end)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(result))
}

func TestDeleteRevenue(t *testing.T) {
	idProviderMock := new(DynamoDbRevenueIdProviderMock)
	repo := newRevenueRepository(idProviderMock)

	input := revenue.Revenue{
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/05/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
	}
	id := revenue.RevenueId("MjAyNF90ZXN0dXNlcg==-NV8xX1NBTFQ=")
	idProviderMock.On("GenerateIdFor", &input).Return(string(id))

	if err := repo.Save(ctx, &input); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	if err := repo.Delete(ctx, input.Id); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	_, err := repo.FindFor(ctx, input.Id)
	assert.NotEqual(t, nil, err)
}

func TestDeleteRevenueFailsWhenNotOwner(t *testing.T) {
	idProviderMock := new(DynamoDbRevenueIdProviderMock)
	repo := newRevenueRepository(idProviderMock)

	input := revenue.Revenue{
		UserName: "testuser",
		Date:     testutils.SafeDateFor("10/05/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "NOTE",
	}
	id := revenue.RevenueId("MjAyNF90ZXN0dXNlcg==-NV8xMF9TQTE=")
	idProviderMock.On("GenerateIdFor", &input).Return(string(id))

	if err := repo.Save(ctx, &input); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	err := repo.Delete(ctxAnotherUser, id)
	assert.NotEqual(t, nil, err)

	retrieved, err := repo.FindFor(ctx, id)
	assert.Equal(t, nil, err)
	assert.Equal(t, "NOTE", retrieved.Note)
}
