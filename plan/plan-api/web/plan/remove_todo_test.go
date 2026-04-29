package plan

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
	"github.com/stretchr/testify/assert"
)

func TestRemoveTodo(t *testing.T) {
	td := &plan.Todo{Id: "test-todo-id", UserName: testUser, Date: clock.ToDay(), Content: "do something"}
	p := aTestPlan()
	p.Todos = []*plan.Todo{td}
	router := setupRouter(&mockRepo{plans: []*plan.Plan{p}})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/plan/"+p.Id+"/todo/"+td.Id, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
