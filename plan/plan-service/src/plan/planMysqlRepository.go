package plan

type MySqlPlanRepository struct {
	ConnectionString string
}

// GetAllPlanBy(userName string) ([]*Plan, error)
func (repostiroy *MySqlPlanRepository) GetPlan(idpalnId string, userName string) (*Plan, error) {
	return nil, nil
}

// GetPlanDetails(id string, userName string) ([]*Todo, error)
func (repostiroy *MySqlPlanRepository) CreateNewPlan(plan Plan) (string, error) {
	return "", nil
}

// AddTodo(id string, todo Todo) error
// RemoveTodo(id string, todoId string) error
