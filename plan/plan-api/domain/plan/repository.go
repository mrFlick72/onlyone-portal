package plan

type PlanRepository interface {
	GetAllPlanBy(userName string) ([]*Plan, error)
	GetPlan(idPlanId string, userName string) (*Plan, error)
	CreateNewPlan(p Plan) (string, error)
	DeletePlan(idPlanId string, userName string) error
	AddTodo(idPlanId string, t Todo) error
	UpdateTodo(idPlanId string, t Todo) error
	RemoveTodo(idPlanId string, todoId string) error
}
