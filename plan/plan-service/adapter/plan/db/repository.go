package db

import (
	domainplan "github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/todo"
)

type PlanPostgresRepository struct {
	ConnectionString string
}

func NewPlanRepository(dsn string) *PlanPostgresRepository {
	return &PlanPostgresRepository{ConnectionString: dsn}
}

func (r *PlanPostgresRepository) GetAllPlanBy(userName string) ([]*domainplan.Plan, error) {
	return nil, nil
}

func (r *PlanPostgresRepository) GetPlan(idPlanId string, userName string) (*domainplan.Plan, error) {
	return nil, nil
}

func (r *PlanPostgresRepository) GetPlanDetails(idPlanId string, userName string) ([]*todo.Todo, error) {
	return nil, nil
}

func (r *PlanPostgresRepository) CreateNewPlan(p domainplan.Plan) (string, error) {
	return "", nil
}

func (r *PlanPostgresRepository) AddTodo(idPlanId string, t todo.Todo) error {
	return nil
}

func (r *PlanPostgresRepository) RemoveTodo(idPlanId string, todoId string) error {
	return nil
}
