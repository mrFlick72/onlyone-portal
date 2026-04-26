package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/pkg/clock"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/plan"
	"github.com/stretchr/testify/assert"
)

const testUserName = "valerio.vaudi"

func setupRouter(mock plan.TodoRepository) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		userName := testUserName
		c.Set("user", security.User{UserName: &userName})
		c.Next()
	})
	RegisterEndpoints(r, &server.GinContextToPlainContextFactory{}, mock)
	return r
}

func TestGetAllTodo(t *testing.T) {
	mock := &mockTodoRepo{todos: []*plan.Todo{aNewTodo()}}
	router := setupRouter(mock)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/todo-service/todo", nil)
	router.ServeHTTP(w, req)

	expected, _ := json.Marshal([]todoRepresentation{fromDomainToRepresentation(mock.todos[0])})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expected), strings.TrimSpace(w.Body.String()))
}

func TestGetAllTodoWhenEmpty(t *testing.T) {
	router := setupRouter(&mockTodoRepo{todos: []*plan.Todo{}})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/todo-service/todo", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", strings.TrimSpace(w.Body.String()))
}

func TestGetOneTodo(t *testing.T) {
	todo := aNewTodo()
	router := setupRouter(&mockTodoRepo{todos: []*plan.Todo{todo}})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/todo-service/todo/"+todo.Id, nil)
	router.ServeHTTP(w, req)

	expected, _ := json.Marshal(fromDomainToRepresentation(todo))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expected), strings.TrimSpace(w.Body.String()))
}

func TestGetOneTodoNotFound(t *testing.T) {
	router := setupRouter(&mockTodoRepo{todos: []*plan.Todo{}})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/todo-service/todo/nonexistent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSaveTodo(t *testing.T) {
	mock := &mockTodoRepo{}
	router := setupRouter(mock)

	body, _ := json.Marshal(todoRepresentation{
		UserName: testUserName,
		Date:     "2026-04-26",
		Content:  "a content",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/todo-service/todo", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestDeleteTodo(t *testing.T) {
	todo := aNewTodo()
	router := setupRouter(&mockTodoRepo{todos: []*plan.Todo{todo}})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/todo-service/todo/"+todo.Id, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func aNewTodo() *plan.Todo {
	id := "test-id-1"
	return &plan.Todo{
		Id:       id,
		UserName: testUserName,
		Date:     clock.ToDay(),
		Content:  "a content",
	}
}

// mockTodoRepo is an in-memory TodoRepository for handler tests.
type mockTodoRepo struct {
	todos []*plan.Todo
}

func (m *mockTodoRepo) GetAllTodo(userName string) ([]*plan.Todo, error) {
	var result []*plan.Todo
	for _, t := range m.todos {
		if t.UserName == userName {
			result = append(result, t)
		}
	}
	if result == nil {
		result = []*plan.Todo{}
	}
	return result, nil
}

func (m *mockTodoRepo) GetTodo(id string) (*plan.Todo, error) {
	for _, t := range m.todos {
		if t.Id == id {
			return t, nil
		}
	}
	return nil, nil
}

func (m *mockTodoRepo) SaveTodo(todo *plan.Todo) error {
	m.todos = append(m.todos, todo)
	return nil
}

func (m *mockTodoRepo) RemoveTodo(id string) error {
	for i, t := range m.todos {
		if t.Id == id {
			m.todos = append(m.todos[:i], m.todos[i+1:]...)
			return nil
		}
	}
	return nil
}

// Ensure mockTodoRepo satisfies the interface at compile time.
var _ plan.TodoRepository = (*mockTodoRepo)(nil)

// Stub context factory (unused in these tests — router middleware sets the user).
type stubContextFactory struct{}

func (s *stubContextFactory) CreateContextFromGin(c *gin.Context) context.Context {
	return server.CopyGinKeysToRequestContext(c)
}
