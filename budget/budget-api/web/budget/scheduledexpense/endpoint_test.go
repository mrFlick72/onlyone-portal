package scheduledexpense

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	domainscheduledexpense "github.com/mrflick72/budget/budget-api/domain/budget/scheduledexpense"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/stretchr/testify/mock"
)

func TestCreateANewScheduledExpense(t *testing.T) {
	r := SetUpRouter()
	facade := new(ScheduledExpenseActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterScheduledExpenseEndpoints(r, contextFactoryConverter, facade)

	domainModel := domainscheduledexpense.ScheduledExpense{
		Description: "Rent",
		Amount:      testutils.SafeMoneyFor("1200.00"),
		Notes:       "Monthly rent",
		Tags:        []tags.SearchTag{},
		Day:         5,
	}

	rep := ScheduledExpenseRepresentation{
		Description: "Rent",
		Amount:      "1200.00",
		Notes:       "Monthly rent",
		Day:         5,
	}
	jsonValue, _ := json.Marshal(rep)

	ctx := testutils.NewStubbedContextWith("USER")
	facade.On("CreateScheduledExpense", ctx, &domainModel).Return(nil)
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)

	req, _ := http.NewRequest("POST", "/api/budget/scheduled-expense", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	facade.AssertCalled(t, "CreateScheduledExpense", ctx, &domainModel)
}

func TestCreateANewScheduledExpenseWithBadAmountReturns400(t *testing.T) {
	r := SetUpRouter()
	facade := new(ScheduledExpenseActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterScheduledExpenseEndpoints(r, contextFactoryConverter, facade)

	rep := ScheduledExpenseRepresentation{
		Description: "Rent",
		Amount:      "not-a-number",
		Day:         5,
	}
	jsonValue, _ := json.Marshal(rep)

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)

	req, _ := http.NewRequest("POST", "/api/budget/scheduled-expense", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	facade.AssertNotCalled(t, "CreateScheduledExpense", mock.Anything, mock.Anything)
}

func TestFindAllScheduledExpenses(t *testing.T) {
	r := SetUpRouter()
	facade := new(ScheduledExpenseActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterScheduledExpenseEndpoints(r, contextFactoryConverter, facade)

	scheduledExpenses := []domainscheduledexpense.ScheduledExpense{
		{
			Id:          "ID1",
			UserName:    "USER",
			Description: "Rent",
			Amount:      testutils.SafeMoneyFor("1200.00"),
			Day:         5,
			Status:      domainscheduledexpense.StatusActive,
		},
	}

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("FindScheduledExpenses", ctx).Return(scheduledExpenses, nil)

	req, _ := http.NewRequest("GET", "/api/budget/scheduled-expense", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var actual ScheduledExpenseListRepresentation
	if err := json.Unmarshal(w.Body.Bytes(), &actual); err != nil {
		t.Fatalf("Error unmarshalling: %v", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, ScheduledExpenseListDomainToRepresentationModel(scheduledExpenses), actual)
	facade.AssertCalled(t, "FindScheduledExpenses", ctx)
}

func TestFindAllScheduledExpensesWhenFacadeFailsReturns500(t *testing.T) {
	r := SetUpRouter()
	facade := new(ScheduledExpenseActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterScheduledExpenseEndpoints(r, contextFactoryConverter, facade)

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("FindScheduledExpenses", ctx).Return(nil, errors.New("find fails"))

	req, _ := http.NewRequest("GET", "/api/budget/scheduled-expense", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
