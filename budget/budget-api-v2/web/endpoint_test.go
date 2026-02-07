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
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}
func TestCreateANewBudgetExpense(t *testing.T) {
	r := SetUpRouter()
	facade := new(BudgetExpenseActionsMock)
	RegisterEndpoints(r, facade)

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
	jsonValue, _ := json.Marshal(budgetExpenseRepresentation)

	ctx := testutils.NewStubbedContextWith("USER")
	facade.On("CreateBudgetExpense", ctx, budgetExpense).Return(nil)

	req, _ := http.NewRequest("POST", "/api/budget/expense", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Body = http.NoBody
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	facade.AssertCalled(t, "CreateBudgetExpense", ctx, budgetExpense)
}
