package plan

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/stretchr/testify/assert"
)

func TestGetOnePlan(t *testing.T) {
	p := aTestPlan()
	repo := &mockRepo{}
	repo.On("GetPlan", p.Id, testUser).Return(p, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/plan/"+p.Id, nil)
	setupRouter(repo).ServeHTTP(w, req)

	expected, _ := json.Marshal(toPlanRepresentation(p))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expected), strings.TrimSpace(w.Body.String()))
	repo.AssertExpectations(t)
}

func TestGetOnePlanNotFound(t *testing.T) {
	repo := &mockRepo{}
	repo.On("GetPlan", "nonexistent", testUser).Return((*plan.Plan)(nil), errors.New("not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/plan/nonexistent", nil)
	setupRouter(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	repo.AssertExpectations(t)
}
