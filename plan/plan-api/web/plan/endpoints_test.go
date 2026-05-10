package plan

import (
	"github.com/gin-gonic/gin"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/stretchr/testify/mock"
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

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetAllPlanBy(userName string) ([]*plan.Plan, error) {
	args := m.Called(userName)
	return args.Get(0).([]*plan.Plan), args.Error(1)
}

func (m *mockRepo) GetPlan(id string, userName string) (*plan.Plan, error) {
	args := m.Called(id, userName)
	return args.Get(0).(*plan.Plan), args.Error(1)
}

func (m *mockRepo) CreateNewPlan(p plan.Plan) (string, error) {
	args := m.Called(p)
	return args.String(0), args.Error(1)
}

func (m *mockRepo) AddTodo(planId string, t plan.Todo) error {
	args := m.Called(planId, t)
	return args.Error(0)
}

func (m *mockRepo) UpdateTodo(planId string, t plan.Todo) error {
	args := m.Called(planId, t)
	return args.Error(0)
}

func (m *mockRepo) UpdateTodoStatus(planId string, todoId string, status plan.TodoStatus) error {
	args := m.Called(planId, todoId, status)
	return args.Error(0)
}

func (m *mockRepo) RemoveTodo(planId string, todoId string) error {
	args := m.Called(planId, todoId)
	return args.Error(0)
}

func (m *mockRepo) DeletePlan(planId string, userName string) error {
	args := m.Called(planId, userName)
	return args.Error(0)
}
