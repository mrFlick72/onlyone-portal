package plan

import (
	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveTodo(t *testing.T) {
	aPlan := test.ANewPlan()
	todoId := "test-todo-id"
	repo := &mockRepo{}
	repo.On("RemoveTodo", aPlan.Id, todoId).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/plan/"+aPlan.Id+"/todo/"+todoId, nil)
	setupRouter(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	repo.AssertExpectations(t)
}
