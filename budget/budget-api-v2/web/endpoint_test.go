package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/stretchr/testify/mock"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestCreateANewBudgetExpense(t *testing.T) {
	r := SetUpRouter()
	facade := new(BudgetExpenseActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterEndpoints(r, contextFactoryConverter, facade)

	budgetExpense := expense.BudgetExpense{
		Date:   testutils.SafeDateFor("01/01/2018"),
		Amount: testutils.SafeMoneyFor("100.00"),
		Note:   "Test note",
		Tag:    "tagKey",
	}

	budgetExpenseRepresentation := BudgetExpenseRepresentation{
		Date:     "01/01/2018",
		Amount:   "100.00",
		Note:     "Test note",
		TagKey:   "tagKey",
		TagValue: "tagValue",
	}
	jsonValue, err := json.Marshal(budgetExpenseRepresentation)
	if err != nil {
		t.Fatalf("Error marshalling JSON: %v", err)
	}

	ctx := testutils.NewStubbedContextWith("USER")
	facade.On("CreateBudgetExpense", ctx, &budgetExpense).Return(nil)
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)

	body := bytes.NewBuffer(jsonValue)
	req, _ := http.NewRequest("POST", "/api/budget/expense", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	facade.AssertCalled(t, "CreateBudgetExpense", ctx, &budgetExpense)
}

func TestUpdateABudgetExpense(t *testing.T) {
	r := SetUpRouter()
	facade := new(BudgetExpenseActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterEndpoints(r, contextFactoryConverter, facade)

	budgetExpense := expense.BudgetExpense{
		Id:     "123-456",
		Date:   testutils.SafeDateFor("01/01/2018"),
		Amount: testutils.SafeMoneyFor("100.00"),
		Note:   "Test note",
		Tag:    "tagKey",
	}

	budgetExpenseRepresentation := BudgetExpenseRepresentation{
		Date:     "01/01/2018",
		Amount:   "100.00",
		Note:     "Test note",
		TagKey:   "tagKey",
		TagValue: "tagValue",
	}
	jsonValue, err := json.Marshal(budgetExpenseRepresentation)
	if err != nil {
		t.Fatalf("Error marshalling JSON: %v", err)
	}

	ctx := testutils.NewStubbedContextWith("USER")
	facade.On("UpdateBudgetExpense", ctx, &budgetExpense).Return(nil)
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)

	body := bytes.NewBuffer(jsonValue)
	req, _ := http.NewRequest("PUT", "/api/budget/expense/123-456", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	facade.AssertCalled(t, "UpdateBudgetExpense", ctx, &budgetExpense)
}

func TestDeleteABudgetExpense(t *testing.T) {
	r := SetUpRouter()
	facade := new(BudgetExpenseActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterEndpoints(r, contextFactoryConverter, facade)

	ctx := testutils.NewStubbedContextWith("USER")
	facade.On("DeleteBudgetExpense", ctx, "123-456").Return(nil)
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)

	req, _ := http.NewRequest("DELETE", "/api/budget/expense/123-456", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
	facade.AssertCalled(t, "DeleteBudgetExpense", ctx, "123-456")
}

func TestFindBudgetExpensesByTimeRange(t *testing.T) {
	r := SetUpRouter()
	facade := new(BudgetExpenseActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterEndpoints(r, contextFactoryConverter, facade)

	expected := &expense.SpentBudget{
		BudgetExpenseList: &[]expense.BudgetExpense{
			{
				Id:       "123-456",
				UserName: "USER",
				Date:     testutils.SafeDateFor("01/01/2018"),
				Amount:   testutils.SafeMoneyFor("100.00"),
				Note:     "Test note",
				Tag:      "tagKey",
			},
			{
				Id:       "123-789",
				UserName: "USER",
				Date:     testutils.SafeDateFor("05/01/2018"),
				Amount:   testutils.SafeMoneyFor("100.00"),
				Note:     "Test note",
				Tag:      "tagKey",
			},
		},

		SearchTags: &map[tags.SearchTagKey]tags.SearchTagValue{
			"tagKey": "tagValue",
		},
	}

	budgetSearchCriteriaRepresentation := BudgetSearchCriteriaRepresentation{
		Year:          "2018",
		Month:         "01",
		SearchTagList: []tags.SearchTagKey{"tagKey"},
	}

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("FindSpentBudget", ctx,
		date.NewMonthFor(budgetSearchCriteriaRepresentation.Month),
		date.NewYearFor(budgetSearchCriteriaRepresentation.Year),
		budgetSearchCriteriaRepresentation.SearchTagList).Return(expected, nil)

	jsonValue, err := json.Marshal(budgetSearchCriteriaRepresentation)
	if err != nil {
		t.Fatalf("Error marshalling JSON: %v", err)
	}

	body := bytes.NewBuffer(jsonValue)
	req, _ := http.NewRequest("PUT", "/api/budget/expense", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var actualBudgetSpent SpentBudgetRepresentation
	err = json.Unmarshal(w.Body.Bytes(), &actualBudgetSpent)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON: %v", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, *SpentBudgetDomainToRepresentationModel(expected), actualBudgetSpent)
	facade.AssertCalled(t, "FindSpentBudget", ctx, date.NewMonthFor(budgetSearchCriteriaRepresentation.Month),
		date.NewYearFor(budgetSearchCriteriaRepresentation.Year),
		budgetSearchCriteriaRepresentation.SearchTagList)
}
