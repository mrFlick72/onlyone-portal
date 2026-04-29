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
	td := &plan.Todo{Id: "test-todo-id", UserName: testUser, Date: clock.ToDay(), Content: "old content"}
	p := aTestPlan()
	p.Todos = []*plan.Todo{td}
	router := setupRouter(&mockRepo{plans: []*plan.Plan{p}})

	body, _ := json.Marshal(todoRepresentation{Date: "2026-04-29", Content: "updated content"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/plan/"+p.Id+"/todo/"+td.Id, strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
