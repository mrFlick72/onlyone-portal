package db

import (
	"database/sql"
	"errors"

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
	result := make([]*domainplan.Plan, 0)

	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return result, err
	}

	query, err := db.Prepare("SELECT id, user_name, title, date FROM plan WHERE user_name = $1")
	if err != nil {
		return result, err
	}

	rows, err := query.Query(userName)
	logger.LogErrorFor(err)
	result = buildPlans(rows)

	database.CloseResources(rows, query, db)

	return result, err
}

func (r *PlanPostgresRepository) GetPlan(idPlanId string, userName string) (*domainplan.Plan, error) {
	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return nil, err
	}

	query, err := db.Prepare("SELECT id, user_name, title, date FROM plan WHERE id = $1 AND user_name = $2")
	if err != nil {
		return nil, err
	}

	rows, err := query.Query(idPlanId, userName)
	logger.LogErrorFor(err)
	result := buildPlans(rows)

	database.CloseResources(rows, query, db)

	if len(result) == 0 {
		return nil, errors.New("plan with id " + idPlanId + " not found")
	}

	plan := result[0]
	plan.Todos, err = r.loadTodosFor(idPlanId)
	return plan, err
}

func (r *PlanPostgresRepository) loadTodosFor(planId string) ([]*todo.Todo, error) {
	result := make([]*todo.Todo, 0)

	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return result, err
	}

	query, err := db.Prepare("SELECT id, user_name, date, content FROM todo WHERE plan_id = $1")
	if err != nil {
		return result, err
	}

	rows, err := query.Query(planId)
	logger.LogErrorFor(err)

	for rows.Next() {
		var t todo.Todo
		rows.Scan(&t.Id, &t.UserName, &t.Date, &t.Content)
		t.Date = t.Date.UTC()
		result = append(result, &t)
	}

	database.CloseResources(rows, query, db)
	return result, err
}

func buildPlans(rows *sql.Rows) []*domainplan.Plan {
	result := make([]*domainplan.Plan, 0)
	for rows.Next() {
		var p domainplan.Plan
		rows.Scan(&p.Id, &p.UserName, &p.Title, &p.Date)
		p.Date = p.Date.UTC()
		p.Todos = []*todo.Todo{}
		result = append(result, &p)
	}
	return result
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
	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return err
	}

	query, err := db.Prepare("INSERT INTO todo (id, plan_id, user_name, date, content) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}

	_, err = query.Exec(t.Id, idPlanId, t.UserName, t.Date, t.Content)
	logger.LogErrorFor(err)
	database.CloseResources(nil, query, db)
	return err
}

func (r *PlanPostgresRepository) RemoveTodo(idPlanId string, todoId string) error {
	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return err
	}

	query, err := db.Prepare("DELETE FROM todo WHERE id = $1 AND plan_id=$2")
	if err != nil {
		return err
	}

	_, err = query.Exec(todoId, idPlanId)
	logger.LogErrorFor(err)
	database.CloseResources(nil, query, db)
	return err
}
