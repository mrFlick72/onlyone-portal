package plan

import "github.com/mrflick72/onlyone-portal/plan/plan-service/domain/todo"

type Repository interface {
	GetAllPlanBy(userName string) ([]*Plan, error)
	GetPlan(idPlanId string, userName string) (*Plan, error)
	GetPlanDetails(idPlanId string, userName string) ([]*todo.Todo, error)
	CreateNewPlan(p Plan) (string, error)
	AddTodo(idPlanId string, t todo.Todo) error
	RemoveTodo(idPlanId string, todoId string) error
}
