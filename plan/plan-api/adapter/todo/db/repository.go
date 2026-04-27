package db

import (
	"errors"

	_ "github.com/lib/pq"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/todo"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/database"
)

var logger = logging.GetLoggerInstanceForComponentByTypeName("TodoPostgresRepository")

type TodoPostgresRepository struct {
	ConnectionString string
}

func NewTodoRepository(dsn string) *TodoPostgresRepository {
	return &TodoPostgresRepository{ConnectionString: dsn}
}

func (r *TodoPostgresRepository) GetAllTodo(userName string) ([]*todo.Todo, error) {
	result := make([]*todo.Todo, 0)

	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return result, err
	}

	query, err := db.Prepare("SELECT id, user_name, date, content FROM todo WHERE user_name = $1")
	if err != nil {
		return result, err
	}

	rows, err := query.Query(userName)
	logger.LogErrorFor(err)
	result = buildTodos(rows, result)

	database.CloseResources(rows, query, db)
	return result, err
}

func (r *TodoPostgresRepository) GetTodo(id string) (*todo.Todo, error) {
	var result []*todo.Todo

	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return nil, err
	}

	query, err := db.Prepare("SELECT id, user_name, date, content FROM todo WHERE id = $1")
	if err != nil {
		return nil, err
	}

	rows, err := query.Query(id)
	logger.LogErrorFor(err)
	result = buildTodos(rows, result)

	database.CloseResources(rows, query, db)

	if len(result) > 0 {
		return result[0], nil
	}
	return nil, errors.New("todo with id " + id + " not found")
}

func (r *TodoPostgresRepository) SaveTodo(t *todo.Todo) error {
	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return err
	}

	query, err := db.Prepare("INSERT INTO todo (id, user_name, date, content) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	_, err = query.Exec(t.Id, t.UserName, t.Date, t.Content)
	logger.LogErrorFor(err)

	database.CloseResources(nil, query, db)
	return err
}

func (r *TodoPostgresRepository) RemoveTodo(id string) error {
	db, err := database.GetDatabaseConnectionFor(r.ConnectionString)
	if err != nil {
		return err
	}

	query, err := db.Prepare("DELETE FROM todo WHERE id = $1")
	if err != nil {
		return err
	}

	_, err = query.Exec(id)
	logger.LogErrorFor(err)

	database.CloseResources(nil, query, db)
	return err
}
