package plan

type TodoRepository interface {
	GetAllTodo(userName string) ([]*Todo, error)
	GetTodo(id string) (*Todo, error)
	SaveTodo(todo *Todo) error
	RemoveTodo(id string) error
}

type PlanRepository interface {
	GetAllPlanBy(userName string) ([]*Plan, error)
	GetPlan(idpalnId string, userName string) (*Plan, error)
	GetPlanDetails(idpalnId string, userName string) ([]*Todo, error)
	CreateNewPlan(plan Plan) (string, error)
	AddTodo(idpalnId string, todo Todo) error
	RemoveTodo(idpalnId string, todoId string) error
}
