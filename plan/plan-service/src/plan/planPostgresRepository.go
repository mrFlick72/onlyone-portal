package plan

type PostgresPlanRepository struct {
	ConnectionString string
}

func (repository *PostgresPlanRepository) GetPlan(idPlanId string, userName string) (*Plan, error) {
	return nil, nil
}

func (repository *PostgresPlanRepository) CreateNewPlan(plan Plan) (string, error) {
	return "", nil
}
