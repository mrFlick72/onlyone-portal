package plan

import (
	"encoding/json"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPlans(t *testing.T) {
	aPlan := test.ANewPlan()
	repo := &mockRepo{}
	repo.On("GetAllPlanBy", testUser).Return([]*plan.Plan{&aPlan}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/plan", nil)
	setupRouter(repo).ServeHTTP(w, req)

	expected, _ := json.Marshal([]planRepresentation{toPlanRepresentation(&aPlan)})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expected), strings.TrimSpace(w.Body.String()))
	repo.AssertExpectations(t)
}

func TestGetAllPlansWhenEmpty(t *testing.T) {
	repo := &mockRepo{}
	repo.On("GetAllPlanBy", testUser).Return([]*plan.Plan{}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/plan", nil)
	setupRouter(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", strings.TrimSpace(w.Body.String()))
	repo.AssertExpectations(t)
}
