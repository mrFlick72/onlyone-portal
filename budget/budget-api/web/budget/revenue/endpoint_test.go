package revenue

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/budget/revenue"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/stretchr/testify/mock"
)

func TestCreateANewRevenue(t *testing.T) {
	r := SetUpRouter()
	facade := new(RevenueActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterRevenueEndpoints(r, contextFactoryConverter, facade)

	domainRevenue := revenue.Revenue{
		Date:   testutils.SafeDateFor("01/01/2018"),
		Amount: testutils.SafeMoneyFor("100.00"),
		Note:   "Test note",
	}

	rep := RevenueRepresentation{
		Date:   "01/01/2018",
		Amount: "100.00",
		Note:   "Test note",
	}
	jsonValue, _ := json.Marshal(rep)

	ctx := testutils.NewStubbedContextWith("USER")
	facade.On("CreateRevenue", ctx, &domainRevenue).Return(nil)
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)

	req, _ := http.NewRequest("POST", "/api/budget/revenue", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	facade.AssertCalled(t, "CreateRevenue", ctx, &domainRevenue)
}

func TestUpdateARevenue(t *testing.T) {
	r := SetUpRouter()
	facade := new(RevenueActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterRevenueEndpoints(r, contextFactoryConverter, facade)

	domainRevenue := revenue.Revenue{
		Id:     "123-456",
		Date:   testutils.SafeDateFor("01/01/2018"),
		Amount: testutils.SafeMoneyFor("100.00"),
		Note:   "Test note",
	}

	rep := RevenueRepresentation{
		Date:   "01/01/2018",
		Amount: "100.00",
		Note:   "Test note",
	}
	jsonValue, _ := json.Marshal(rep)

	ctx := testutils.NewStubbedContextWith("USER")
	facade.On("UpdateRevenue", ctx, &domainRevenue).Return(nil)
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)

	req, _ := http.NewRequest("PUT", "/api/budget/revenue/123-456", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	facade.AssertCalled(t, "UpdateRevenue", ctx, &domainRevenue)
}

func TestDeleteARevenue(t *testing.T) {
	r := SetUpRouter()
	facade := new(RevenueActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterRevenueEndpoints(r, contextFactoryConverter, facade)

	ctx := testutils.NewStubbedContextWith("USER")
	facade.On("DeleteRevenue", ctx, "123-456").Return(nil)
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)

	req, _ := http.NewRequest("DELETE", "/api/budget/revenue/123-456", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	facade.AssertCalled(t, "DeleteRevenue", ctx, "123-456")
}

func TestFindRevenuesByYear(t *testing.T) {
	r := SetUpRouter()
	facade := new(RevenueActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterRevenueEndpoints(r, contextFactoryConverter, facade)

	expected := []revenue.Revenue{
		{
			Id:       "123-456",
			UserName: "USER",
			Date:     testutils.SafeDateFor("01/01/2018"),
			Amount:   testutils.SafeMoneyFor("100.00"),
			Note:     "Test note",
		},
	}

	ctx := testutils.NewStubbedContextWith("USER")
	contextFactoryConverter.On("CreateContextFromGin", mock.AnythingOfType("*gin.Context")).Return(ctx)
	facade.On("FindRevenue", ctx, date.NewYear(2018)).Return(expected, nil)

	req, _ := http.NewRequest("GET", "/api/budget/revenue?q=year=2018", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var actual []RevenueRepresentation
	if err := json.Unmarshal(w.Body.Bytes(), &actual); err != nil {
		t.Fatalf("Error unmarshalling: %v", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, RevenueListDomainToRepresentationModel(expected), actual)
	facade.AssertCalled(t, "FindRevenue", ctx, date.NewYear(2018))
}

func TestFindRevenuesByYearBadQueryReturns400(t *testing.T) {
	r := SetUpRouter()
	facade := new(RevenueActionsMock)
	contextFactoryConverter := new(ContextFactoryConverterMock)
	RegisterRevenueEndpoints(r, contextFactoryConverter, facade)

	req, _ := http.NewRequest("GET", "/api/budget/revenue", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	facade.AssertNotCalled(t, "FindRevenue", mock.Anything, mock.Anything)
}
