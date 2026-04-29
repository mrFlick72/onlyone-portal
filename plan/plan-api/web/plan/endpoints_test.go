package plan

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
	"github.com/stretchr/testify/assert"
)

const testUser = "valerio.vaudi"

func setupRouter(repo plan.PlanRepository) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		userName := testUser
		c.Set("user", security.User{UserName: &userName})
		c.Next()
	})
	RegisterEndpoints(r, &server.GinContextToPlainContextFactory{}, repo)
	return r
}

func aTestPlan() *plan.Plan {
	return &plan.Plan{
		Id:       "test-plan-id",
		UserName: testUser,
		Title:    "test plan",
		Date:     clock.ToDay(),
		Todos:    []*plan.Todo{},
	}
}

// --- GET /api/plan ---

func TestGetAllPlans(t *testing.T) {
	p := aTestPlan()
	router := setupRouter(&mockRepo{plans: []*plan.Plan{p}})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/plan", nil)
	router.ServeHTTP(w, req)

	expected, _ := json.Marshal([]planRepresentation{toPlanRepresentation(p)})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expected), strings.TrimSpace(w.Body.String()))
}

func TestGetAllPlansWhenEmpty(t *testing.T) {
	router := setupRouter(&mockRepo{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/plan", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", strings.TrimSpace(w.Body.String()))
}

// --- GET /api/plan/:id ---

func TestGetOnePlan(t *testing.T) {
	p := aTestPlan()
	router := setupRouter(&mockRepo{plans: []*plan.Plan{p}})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/plan/"+p.Id, nil)
	router.ServeHTTP(w, req)

	expected, _ := json.Marshal(toPlanRepresentation(p))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expected), strings.TrimSpace(w.Body.String()))
}

func TestGetOnePlanNotFound(t *testing.T) {
	router := setupRouter(&mockRepo{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/plan/nonexistent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// --- POST /api/plan ---

func TestSavePlan(t *testing.T) {
	router := setupRouter(&mockRepo{})

	body, _ := json.Marshal(planRepresentation{Title: "new plan", Date: "2026-04-29"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/plan", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	_, err := uuid.Parse(resp["id"])
	assert.NoError(t, err, "response id should be a valid UUID")
}

// --- PUT /api/plan/:id/todo ---

func TestAddTodo(t *testing.T) {
	p := aTestPlan()
	router := setupRouter(&mockRepo{plans: []*plan.Plan{p}})

	body, _ := json.Marshal(todoRepresentation{Date: "2026-04-29", Content: "do something"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/plan/"+p.Id+"/todo", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// --- DELETE /api/plan/:id/todo/:todoId ---

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

// --- PUT /api/plan/:id/todo/:todoId ---

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

// --- DELETE /api/plan/:id ---

func TestDeletePlan(t *testing.T) {
	p := aTestPlan()
	router := setupRouter(&mockRepo{plans: []*plan.Plan{p}})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/plan/"+p.Id, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// --- mock ---

type mockRepo struct {
	plans []*plan.Plan
}

func (m *mockRepo) GetAllPlanBy(userName string) ([]*plan.Plan, error) {
	result := make([]*plan.Plan, 0)
	for _, p := range m.plans {
		if p.UserName == userName {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockRepo) GetPlan(id string, userName string) (*plan.Plan, error) {
	for _, p := range m.plans {
		if p.Id == id && p.UserName == userName {
			return p, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockRepo) CreateNewPlan(p plan.Plan) (string, error) {
	id := uuid.NewString()
	p.Id = id
	p.Todos = []*plan.Todo{}
	m.plans = append(m.plans, &p)
	return id, nil
}

func (m *mockRepo) AddTodo(planId string, t plan.Todo) error {
	for _, p := range m.plans {
		if p.Id == planId {
			p.Todos = append(p.Todos, &t)
			return nil
		}
	}
	return errors.New("plan not found")
}

func (m *mockRepo) UpdateTodo(planId string, t plan.Todo) error {
	for _, p := range m.plans {
		if p.Id == planId {
			for i, td := range p.Todos {
				if td.Id == t.Id {
					p.Todos[i] = &t
					return nil
				}
			}
		}
	}
	return errors.New("todo not found")
}

func (m *mockRepo) RemoveTodo(planId string, todoId string) error {
	for _, p := range m.plans {
		if p.Id == planId {
			for i, t := range p.Todos {
				if t.Id == todoId {
					p.Todos = append(p.Todos[:i], p.Todos[i+1:]...)
					return nil
				}
			}
		}
	}
	return nil
}

func (m *mockRepo) DeletePlan(planId string, userName string) error {
	for i, p := range m.plans {
		if p.Id == planId && p.UserName == userName {
			m.plans = append(m.plans[:i], m.plans[i+1:]...)
			return nil
		}
	}
	return errors.New("plan not found")
}

var _ plan.PlanRepository = (*mockRepo)(nil)
