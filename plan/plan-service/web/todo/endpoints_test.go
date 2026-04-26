package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/todo"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
	"github.com/stretchr/testify/assert"
)

const testUser = "valerio.vaudi"

func setupRouter(repo todo.Repository) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		userName := testUser
		c.Set("user", security.User{UserName: &userName})
		c.Next()
	})
	RegisterEndpoints(r, &server.GinContextToPlainContextFactory{}, repo)
	return r
}

func TestGetAllTodo(t *testing.T) {
	td := aTestTodo()
	router := setupRouter(&mockRepo{todos: []*todo.Todo{td}})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/todo-service/todo", nil)
	router.ServeHTTP(w, req)

	expected, _ := json.Marshal([]todoRepresentation{toRepresentation(td)})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expected), strings.TrimSpace(w.Body.String()))
}

func TestGetAllTodoWhenEmpty(t *testing.T) {
	router := setupRouter(&mockRepo{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/todo-service/todo", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", strings.TrimSpace(w.Body.String()))
}

func TestGetOneTodo(t *testing.T) {
	td := aTestTodo()
	router := setupRouter(&mockRepo{todos: []*todo.Todo{td}})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/todo-service/todo/"+td.Id, nil)
	router.ServeHTTP(w, req)

	expected, _ := json.Marshal(toRepresentation(td))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expected), strings.TrimSpace(w.Body.String()))
}

func TestGetOneTodoNotFound(t *testing.T) {
	router := setupRouter(&mockRepo{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/todo-service/todo/nonexistent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSaveTodo(t *testing.T) {
	router := setupRouter(&mockRepo{})

	body, _ := json.Marshal(todoRepresentation{UserName: testUser, Date: "2026-04-26", Content: "a content"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/todo-service/todo", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestDeleteTodo(t *testing.T) {
	td := aTestTodo()
	router := setupRouter(&mockRepo{todos: []*todo.Todo{td}})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/todo-service/todo/"+td.Id, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func aTestTodo() *todo.Todo {
	return &todo.Todo{Id: "test-id-1", UserName: testUser, Date: clock.ToDay(), Content: "a content"}
}

type mockRepo struct {
	todos []*todo.Todo
}

func (m *mockRepo) GetAllTodo(userName string) ([]*todo.Todo, error) {
	var result []*todo.Todo
	for _, t := range m.todos {
		if t.UserName == userName {
			result = append(result, t)
		}
	}
	if result == nil {
		result = []*todo.Todo{}
	}
	return result, nil
}

func (m *mockRepo) GetTodo(id string) (*todo.Todo, error) {
	for _, t := range m.todos {
		if t.Id == id {
			return t, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) SaveTodo(t *todo.Todo) error {
	m.todos = append(m.todos, t)
	return nil
}

func (m *mockRepo) RemoveTodo(id string) error {
	for i, t := range m.todos {
		if t.Id == id {
			m.todos = append(m.todos[:i], m.todos[i+1:]...)
			return nil
		}
	}
	return nil
}

var _ todo.Repository = (*mockRepo)(nil)
