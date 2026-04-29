package plan

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSavePlan(t *testing.T) {
	newId := uuid.NewString()
	repo := &mockRepo{}
	repo.On("CreateNewPlan", mock.MatchedBy(func(p plan.Plan) bool {
		return p.UserName == testUser && p.Title == "new plan" && p.Date == clock.ParseDateFor("2026-04-29")
	})).Return(newId, nil)

	body, _ := json.Marshal(planRepresentation{Title: "new plan", Date: "2026-04-29"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/plan", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, newId, resp["id"])
	repo.AssertExpectations(t)
}
