package plan

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/stretchr/testify/assert"
)

func TestGetOnePlan(t *testing.T) {
	p := aTestPlan()
	router := setupRouter(&mockRepo{plans: []*plan.Plan{p}})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/plan/"+p.Id, nil)
	router.ServeHTTP(w, req)

	expected, _ := json.Marshal(toPlanRepresentation(p))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expected), strings.TrimSpace(w.Body.String()))
}

func TestGetOnePlanNotFound(t *testing.T) {
	router := setupRouter(&mockRepo{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/plan/nonexistent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
