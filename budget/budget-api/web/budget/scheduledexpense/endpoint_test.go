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

// TestCreateANewScheduledExpenseWhenFacadeFailsReturns500 guards against the
// handler discarding Execute's error and reporting success regardless — the
// exact bug a misconfigured/missing DynamoDB table name triggers (Save fails,
// the client is told 201 anyway). See web/budget/expense and
// web/budget/revenue's Create/Update/Delete handlers for the same
// pre-existing shape; not fixed here (out of #51's scope).
func TestCreateANewScheduledExpenseWhenFacadeFailsReturns500(t *testing.T) {
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
	facade.On("CreateScheduledExpense", ctx, &domainModel).Return(errors.New("ResourceNotFoundException: table not found"))
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)

	req, _ := http.NewRequest("POST", "/api/budget/scheduled-expense", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateAScheduledExpense(t *testing.T) {
	r := SetUpRouter()
	facade := new(ScheduledExpenseActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterScheduledExpenseEndpoints(r, contextFactoryConverter, facade)

	domainModel := domainscheduledexpense.ScheduledExpense{
		Id:          "123-456",
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
	facade.On("UpdateScheduledExpense", ctx, &domainModel).Return(nil)
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)

	req, _ := http.NewRequest("PUT", "/api/budget/scheduled-expense/123-456", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	facade.AssertCalled(t, "UpdateScheduledExpense", ctx, &domainModel)
}

func TestUpdateAScheduledExpenseWithBadAmountReturns400(t *testing.T) {
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

	req, _ := http.NewRequest("PUT", "/api/budget/scheduled-expense/123-456", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	facade.AssertNotCalled(t, "UpdateScheduledExpense", mock.Anything, mock.Anything)
}

func TestUpdateAScheduledExpenseWhenFacadeFailsReturns500(t *testing.T) {
	r := SetUpRouter()
	facade := new(ScheduledExpenseActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterScheduledExpenseEndpoints(r, contextFactoryConverter, facade)

	domainModel := domainscheduledexpense.ScheduledExpense{
		Id:          "123-456",
		Description: "Rent",
		Amount:      testutils.SafeMoneyFor("1200.00"),
		Day:         5,
		Tags:        []tags.SearchTag{},
	}

	rep := ScheduledExpenseRepresentation{
		Description: "Rent",
		Amount:      "1200.00",
		Day:         5,
	}
	jsonValue, _ := json.Marshal(rep)

	ctx := testutils.NewStubbedContextWith("USER")
	facade.On("UpdateScheduledExpense", ctx, &domainModel).Return(errors.New("not found or not authorized"))
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)

	req, _ := http.NewRequest("PUT", "/api/budget/scheduled-expense/123-456", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetScheduledExpenseById(t *testing.T) {
	r := SetUpRouter()
	facade := new(ScheduledExpenseActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterScheduledExpenseEndpoints(r, contextFactoryConverter, facade)

	scheduledExpense := &domainscheduledexpense.ScheduledExpense{
		Id:          "123-456",
		UserName:    "USER",
		Description: "Rent",
		Amount:      testutils.SafeMoneyFor("1200.00"),
		Day:         5,
		Status:      domainscheduledexpense.StatusActive,
	}

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("FindScheduledExpense", ctx, domainscheduledexpense.ScheduledExpenseId("123-456")).Return(scheduledExpense, nil)

	req, _ := http.NewRequest("GET", "/api/budget/scheduled-expense/123-456", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var actual ScheduledExpenseRepresentation
	if err := json.Unmarshal(w.Body.Bytes(), &actual); err != nil {
		t.Fatalf("Error unmarshalling: %v", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, ScheduledExpenseDomainToRepresentationModel(scheduledExpense), actual)
}

func TestGetScheduledExpenseByIdWhenNotFoundReturns404(t *testing.T) {
	r := SetUpRouter()
	facade := new(ScheduledExpenseActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterScheduledExpenseEndpoints(r, contextFactoryConverter, facade)

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("FindScheduledExpense", ctx, domainscheduledexpense.ScheduledExpenseId("MISSING")).Return(nil, nil)

	req, _ := http.NewRequest("GET", "/api/budget/scheduled-expense/MISSING", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetScheduledExpenseByIdWhenFacadeFailsReturns500(t *testing.T) {
	r := SetUpRouter()
	facade := new(ScheduledExpenseActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterScheduledExpenseEndpoints(r, contextFactoryConverter, facade)

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("FindScheduledExpense", ctx, domainscheduledexpense.ScheduledExpenseId("123-456")).Return(nil, errors.New("find fails"))

	req, _ := http.NewRequest("GET", "/api/budget/scheduled-expense/123-456", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
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
