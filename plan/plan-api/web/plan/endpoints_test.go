package plan

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
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
