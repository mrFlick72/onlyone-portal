package plan

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveTodo(t *testing.T) {
	p := aTestPlan()
	todoId := "test-todo-id"
	repo := &mockRepo{}
	repo.On("RemoveTodo", p.Id, todoId).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/plan/"+p.Id+"/todo/"+todoId, nil)
	setupRouter(repo).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	repo.AssertExpectations(t)
}
