package db

import (
	"database/sql"
	"errors"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"

	"github.com/google/uuid"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/database"
)

var logger = logging.GetLoggerInstanceForComponentByTypeName("PlanPostgresRepository")

type PlanPostgresRepository struct {
	ConnectionString string
}

func (r *PlanPostgresRepository) GetAllPlanBy(userName string) ([]*plan.Plan, error) {
	result := make([]*plan.Plan, 0)

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

	for _, p := range result {
		p.Todos, err = r.loadTodosFor(p.Id)
		if err != nil {
			return result, err
		}
	}

	return result, err
}

func (r *PlanPostgresRepository) GetPlan(idPlanId string, userName string) (*plan.Plan, error) {
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
	results := buildPlans(rows)

	database.CloseResources(rows, query, db)

	if len(results) == 0 {
		return nil, errors.New("plan with id " + idPlanId + " not found")
	}

	result := results[0]
	result.Todos, err = r.loadTodosFor(idPlanId)
	return result, err
}

func (r *PlanPostgresRepository) loadTodosFor(planId string) ([]*plan.Todo, error) {
	result := make([]*plan.Todo, 0)

	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return result, err
	}

	query, err := db.Prepare("SELECT id, user_name, date, content, status FROM todo WHERE plan_id = $1")
	if err != nil {
		return result, err
	}

	rows, err := query.Query(planId)
	logger.LogErrorFor(err)

	for rows.Next() {
		var t plan.Todo
		var status string
		//todo han\dle the error there
		rows.Scan(&t.Id, &t.UserName, &t.Date, &t.Content, &status)
		t.Date = t.Date.UTC()
		t.Status = plan.TodoStatus(status)
		result = append(result, &t)
	}

	database.CloseResources(rows, query, db)
	return result, err
}

func buildPlans(rows *sql.Rows) []*plan.Plan {
	result := make([]*plan.Plan, 0)
	for rows.Next() {
		var p plan.Plan
		rows.Scan(&p.Id, &p.UserName, &p.Title, &p.Date)
		p.Date = p.Date.UTC()
		p.Todos = []*plan.Todo{}
		result = append(result, &p)
	}
	return result
}

func (r *PlanPostgresRepository) CreateNewPlan(p plan.Plan) (string, error) {
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

func (r *PlanPostgresRepository) DeletePlan(idPlanId string, userName string) error {
	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return err
	}

	query, err := db.Prepare("DELETE FROM plan WHERE id = $1 AND user_name = $2")
	if err != nil {
		return err
	}

	_, err = query.Exec(idPlanId, userName)
	logger.LogErrorFor(err)
	database.CloseResources(nil, query, db)
	return err
}

func (r *PlanPostgresRepository) AddTodo(idPlanId string, t plan.Todo) error {
	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return err
	}

	query, err := db.Prepare("INSERT INTO todo (id, plan_id, user_name, date, content, status) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return err
	}

	_, err = query.Exec(t.Id, idPlanId, t.UserName, t.Date, t.Content, string(t.Status))
	logger.LogErrorFor(err)
	database.CloseResources(nil, query, db)
	return err
}

func (r *PlanPostgresRepository) UpdateTodo(idPlanId string, t plan.Todo) error {
	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return err
	}

	query, err := db.Prepare("UPDATE todo SET content = $1, date = $2 WHERE id = $3 AND plan_id = $4")
	if err != nil {
		return err
	}

	_, err = query.Exec(t.Content, t.Date, t.Id, idPlanId)
	logger.LogErrorFor(err)
	database.CloseResources(nil, query, db)
	return err
}

func (r *PlanPostgresRepository) UpdateTodoStatus(idPlanId string, todoId string, status plan.TodoStatus) error {
	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return err
	}

	query, err := db.Prepare("UPDATE todo SET status = $1 WHERE id = $2 AND plan_id = $3")
	if err != nil {
		return err
	}

	_, err = query.Exec(string(status), todoId, idPlanId)
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
