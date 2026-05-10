package plan

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
	"github.com/stretchr/testify/assert"
)

func aPlanWithTodo(todoId string, status plan.TodoStatus) *plan.Plan {
	p := test.ANewPlan()
	return &plan.Plan{
		Id:       p.Id,
		UserName: testUser,
		Title:    p.Title,
		Date:     p.Date,
		Todos: []*plan.Todo{
			{Id: todoId, UserName: testUser, Date: p.Date, Content: "x", Status: status},
		},
	}
}

func performStatusChange(t *testing.T, repo plan.PlanRepository, planId, todoId, status string) *httptest.ResponseRecorder {
	body, _ := json.Marshal(todoStatusChangeRequest{Status: status})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/plan/"+planId+"/todo/"+todoId+"/status", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(repo).ServeHTTP(w, req)
	return w
}

func TestChangeTodoStatusFromTodoToInProgress(t *testing.T) {
	todoId := "t-1"
	p := aPlanWithTodo(todoId, plan.StatusTodo)
	repo := &mockRepo{}
	repo.On("GetPlan", p.Id, testUser).Return(p, nil)
	repo.On("UpdateTodoStatus", p.Id, todoId, plan.StatusInProgress).Return(nil)

	w := performStatusChange(t, repo, p.Id, todoId, "IN_PROGRESS")

	assert.Equal(t, http.StatusNoContent, w.Code)
	repo.AssertExpectations(t)
}

func TestChangeTodoStatusFromInProgressToDone(t *testing.T) {
	todoId := "t-2"
	p := aPlanWithTodo(todoId, plan.StatusInProgress)
	repo := &mockRepo{}
	repo.On("GetPlan", p.Id, testUser).Return(p, nil)
	repo.On("UpdateTodoStatus", p.Id, todoId, plan.StatusDone).Return(nil)

	w := performStatusChange(t, repo, p.Id, todoId, "DONE")

	assert.Equal(t, http.StatusNoContent, w.Code)
	repo.AssertExpectations(t)
}

func TestChangeTodoStatusToAbortedFromTodo(t *testing.T) {
	todoId := "t-3"
	p := aPlanWithTodo(todoId, plan.StatusTodo)
	repo := &mockRepo{}
	repo.On("GetPlan", p.Id, testUser).Return(p, nil)
	repo.On("UpdateTodoStatus", p.Id, todoId, plan.StatusAborted).Return(nil)

	w := performStatusChange(t, repo, p.Id, todoId, "ABORTED")

	assert.Equal(t, http.StatusNoContent, w.Code)
	repo.AssertExpectations(t)
}

func TestChangeTodoStatusRejectsInvalidTransition(t *testing.T) {
	todoId := "t-4"
	p := aPlanWithTodo(todoId, plan.StatusTodo)
	repo := &mockRepo{}
	repo.On("GetPlan", p.Id, testUser).Return(p, nil)

	w := performStatusChange(t, repo, p.Id, todoId, "DONE")

	assert.Equal(t, http.StatusConflict, w.Code)
	repo.AssertNotCalled(t, "UpdateTodoStatus")
}

func TestChangeTodoStatusRejectsTransitionFromTerminalState(t *testing.T) {
	todoId := "t-5"
	p := aPlanWithTodo(todoId, plan.StatusDone)
	repo := &mockRepo{}
	repo.On("GetPlan", p.Id, testUser).Return(p, nil)

	w := performStatusChange(t, repo, p.Id, todoId, "IN_PROGRESS")

	assert.Equal(t, http.StatusConflict, w.Code)
	repo.AssertNotCalled(t, "UpdateTodoStatus")
}

func TestChangeTodoStatusRejectsUnknownStatus(t *testing.T) {
	todoId := "t-6"
	p := aPlanWithTodo(todoId, plan.StatusTodo)
	repo := &mockRepo{}

	w := performStatusChange(t, repo, p.Id, todoId, "PENDING")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	repo.AssertNotCalled(t, "GetPlan")
	repo.AssertNotCalled(t, "UpdateTodoStatus")
}

func TestChangeTodoStatusReturnsNotFoundWhenPlanMissing(t *testing.T) {
	repo := &mockRepo{}
	repo.On("GetPlan", "missing", testUser).Return((*plan.Plan)(nil), errors.New("not found"))

	w := performStatusChange(t, repo, "missing", "any", "IN_PROGRESS")

	assert.Equal(t, http.StatusNotFound, w.Code)
	repo.AssertNotCalled(t, "UpdateTodoStatus")
}

func TestChangeTodoStatusReturnsNotFoundWhenTodoMissing(t *testing.T) {
	p := aPlanWithTodo("other-todo", plan.StatusTodo)
	repo := &mockRepo{}
	repo.On("GetPlan", p.Id, testUser).Return(p, nil)

	w := performStatusChange(t, repo, p.Id, "missing-todo", "IN_PROGRESS")

	assert.Equal(t, http.StatusNotFound, w.Code)
	repo.AssertNotCalled(t, "UpdateTodoStatus")
}
