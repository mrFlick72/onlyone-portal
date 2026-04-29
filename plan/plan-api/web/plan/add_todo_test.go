package plan

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddTodo(t *testing.T) {
	p := aTestPlan()
	repo := &mockRepo{}
	repo.On("AddTodo", p.Id, mock.MatchedBy(func(t plan.Todo) bool {
		return t.UserName == testUser && t.Content == "do something" && t.Date == clock.ParseDateFor("2026-04-29")
	})).Return(nil)

	body, _ := json.Marshal(todoRepresentation{Date: "2026-04-29", Content: "do something"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/plan/"+p.Id+"/todo", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	repo.AssertExpectations(t)
}
