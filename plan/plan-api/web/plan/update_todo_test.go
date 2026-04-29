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
)

func TestUpdateTodo(t *testing.T) {
	p := aTestPlan()
	todoId := "test-todo-id"
	repo := &mockRepo{}
	repo.On("UpdateTodo", p.Id, plan.Todo{
		Id:       todoId,
		UserName: testUser,
		Date:     clock.ParseDateFor("2026-04-29"),
		Content:  "updated content",
	}).Return(nil)

	body, _ := json.Marshal(todoRepresentation{Date: "2026-04-29", Content: "updated content"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/plan/"+p.Id+"/todo/"+todoId, strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	repo.AssertExpectations(t)
}
