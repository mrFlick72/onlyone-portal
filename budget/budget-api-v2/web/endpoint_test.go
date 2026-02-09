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
