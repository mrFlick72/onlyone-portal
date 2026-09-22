package scheduledexpense

import (
	"testing"

	"github.com/go-playground/assert/v2"
	domainscheduledexpense "github.com/mrflick72/budget/budget-api/domain/budget/scheduledexpense"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	tagRep "github.com/mrflick72/budget/budget-api/web/tags"
)

func TestScheduledExpenseRepresentationToDomainModel(t *testing.T) {
	rep := ScheduledExpenseRepresentation{
		Description: "Rent",
		Amount:      "1200.00",
		Notes:       "Monthly rent",
		Tags:        []tagRep.SearchTagRepresentation{{Key: "housing", Value: "Housing"}},
		Day:         5,
	}

	domainModel, err := ScheduledExpenseRepresentationToDomainModel(rep)

	assert.Equal(t, nil, err)
	assert.Equal(t, "Rent", domainModel.Description)
	assert.Equal(t, testutils.SafeMoneyFor("1200.00"), domainModel.Amount)
	assert.Equal(t, "Monthly rent", domainModel.Notes)
	assert.Equal(t, []tags.SearchTag{{Key: "housing", Value: "Housing"}}, domainModel.Tags)
	assert.Equal(t, 5, domainModel.Day)
	assert.Equal(t, true, domainModel.Month == nil)
	assert.Equal(t, true, domainModel.EndDate == nil)
}

func TestScheduledExpenseRepresentationToDomainModelWithOptionalFields(t *testing.T) {
	month := 3
	endDate := "31/12/2026"
	rep := ScheduledExpenseRepresentation{
		Description: "Insurance",
		Amount:      "450.00",
		Day:         15,
		Month:       &month,
		EndDate:     &endDate,
	}

	domainModel, err := ScheduledExpenseRepresentationToDomainModel(rep)

	assert.Equal(t, nil, err)
	if domainModel.Month == nil || *domainModel.Month != month {
		t.Fatalf("Expected Month %d, got %v", month, domainModel.Month)
	}
	if domainModel.EndDate == nil || domainModel.EndDate.GetFormattedDate() != endDate {
		t.Fatalf("Expected EndDate %v, got %v", endDate, domainModel.EndDate)
	}
}

func TestScheduledExpenseRepresentationToDomainModelWithBadEndDateReturnsError(t *testing.T) {
	badEndDate := "not-a-date"
	rep := ScheduledExpenseRepresentation{
		Description: "Insurance",
		Amount:      "450.00",
		Day:         15,
		EndDate:     &badEndDate,
	}

	_, err := ScheduledExpenseRepresentationToDomainModel(rep)

	assert.NotEqual(t, nil, err)
}

func TestScheduledExpenseRepresentationToDomainModelWithBadAmountReturnsError(t *testing.T) {
	rep := ScheduledExpenseRepresentation{
		Description: "Rent",
		Amount:      "not-a-number",
		Day:         5,
	}

	_, err := ScheduledExpenseRepresentationToDomainModel(rep)

	assert.NotEqual(t, nil, err)
}

func TestScheduledExpenseDomainToRepresentationModel(t *testing.T) {
	domainModel := &domainscheduledexpense.ScheduledExpense{
		Id:          "ID1",
		Description: "Rent",
		Amount:      testutils.SafeMoneyFor("1200.00"),
		Notes:       "Monthly rent",
		Tags:        []tags.SearchTag{{Key: "housing", Value: "Housing"}},
		Day:         5,
		Status:      domainscheduledexpense.StatusActive,
	}

	rep := ScheduledExpenseDomainToRepresentationModel(domainModel)

	assert.Equal(t, "ID1", rep.Id)
	assert.Equal(t, "Rent", rep.Description)
	assert.Equal(t, "1200.00", rep.Amount)
	assert.Equal(t, "Monthly rent", rep.Notes)
	assert.Equal(t, []tagRep.SearchTagRepresentation{{Key: "housing", Value: "Housing"}}, rep.Tags)
	assert.Equal(t, 5, rep.Day)
	assert.Equal(t, "ACTIVE", rep.Status)
	assert.Equal(t, true, rep.Month == nil)
	assert.Equal(t, true, rep.EndDate == nil)
}

func TestScheduledExpenseDomainToRepresentationModelWithOptionalFields(t *testing.T) {
	month := 3
	endDate := testutils.SafeDateFor("31/12/2026")
	domainModel := &domainscheduledexpense.ScheduledExpense{
		Id:          "ID2",
		Description: "Insurance",
		Amount:      testutils.SafeMoneyFor("450.00"),
		Day:         15,
		Month:       &month,
		EndDate:     &endDate,
		Status:      domainscheduledexpense.StatusPaused,
	}

	rep := ScheduledExpenseDomainToRepresentationModel(domainModel)

	if rep.Month == nil || *rep.Month != month {
		t.Fatalf("Expected Month %d, got %v", month, rep.Month)
	}
	if rep.EndDate == nil || *rep.EndDate != "31/12/2026" {
		t.Fatalf("Expected EndDate 31/12/2026, got %v", rep.EndDate)
	}
	assert.Equal(t, "PAUSED", rep.Status)
}

func TestScheduledExpenseListDomainToRepresentationModel(t *testing.T) {
	list := []domainscheduledexpense.ScheduledExpense{
		{Id: "ID1", Description: "Rent", Amount: testutils.SafeMoneyFor("1200.00"), Day: 5, Status: domainscheduledexpense.StatusActive},
		{Id: "ID2", Description: "Insurance", Amount: testutils.SafeMoneyFor("450.00"), Day: 15, Status: domainscheduledexpense.StatusActive},
	}

	rep := ScheduledExpenseListDomainToRepresentationModel(list)

	assert.Equal(t, 2, len(rep.ScheduledExpenses))
	assert.Equal(t, "ID1", rep.ScheduledExpenses[0].Id)
	assert.Equal(t, "ID2", rep.ScheduledExpenses[1].Id)
}
