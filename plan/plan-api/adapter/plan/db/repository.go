package db

import (
	"github.com/google/uuid"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	domainplan "github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/todo"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/database"
)

var logger = logging.GetLoggerInstanceForComponentByTypeName("PlanPostgresRepository")

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
	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return "", err
	}

	id, _ := uuid.NewRandom()
	planId := id.String()

	query, err := db.Prepare("INSERT INTO plan (id, user_name, title, date) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return "", err
	}

	_, err = query.Exec(planId, p.UserName, p.Title, p.Date)
	logger.LogErrorFor(err)
	database.CloseResources(nil, query, db)
	return planId, err
}

func (r *PlanPostgresRepository) AddTodo(idPlanId string, t todo.Todo) error {
	return nil
}

func (r *PlanPostgresRepository) RemoveTodo(idPlanId string, todoId string) error {
	return nil
}
