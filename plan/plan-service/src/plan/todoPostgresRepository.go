package plan

import (
	"errors"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/pkg/database"
)

type PostgresTodoRepository struct {
	ConnectionString string
}

func NewPostgresTodoRepository(dsn string) *PostgresTodoRepository {
	return &PostgresTodoRepository{ConnectionString: dsn}
}

func (repository *PostgresTodoRepository) GetAllTodo(userName string) ([]*Todo, error) {
	result := make([]*Todo, 0)

	db, err := database.GetDatabaseConnectionFor(repository.ConnectionString)
	if err != nil {
		return result, err
	}

	query, err := db.Prepare("SELECT id, user_name, date, content FROM todo WHERE user_name = $1")
	if err != nil {
		return result, err
	}

	rows, err := query.Query(userName)
	logger.LogErrorFor(err)
	result = BuildTodos(rows, result)

	database.CloseResources(rows, query, db)
	return result, err
}

func (repository *PostgresTodoRepository) GetTodo(id string) (*Todo, error) {
	var result []*Todo

	db, err := database.GetDatabaseConnectionFor(repository.ConnectionString)
	if err != nil {
		return nil, err
	}

	query, err := db.Prepare("SELECT id, user_name, date, content FROM todo WHERE id = $1")
	if err != nil {
		return nil, err
	}

	rows, err := query.Query(id)
	logger.LogErrorFor(err)
	result = BuildTodos(rows, result)

	database.CloseResources(rows, query, db)

	if len(result) > 0 {
		return result[0], nil
	}
	return nil, errors.New("todo with id " + id + " not found")
}

func (repository *PostgresTodoRepository) SaveTodo(todo *Todo) error {
	db, err := database.GetDatabaseConnectionFor(repository.ConnectionString)
	if err != nil {
		return err
	}

	query, err := db.Prepare("INSERT INTO todo (id, user_name, date, content) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	_, err = query.Exec(todo.Id, todo.UserName, todo.Date, todo.Content)
	logger.LogErrorFor(err)

	database.CloseResources(nil, query, db)
	return err
}

func (repository *PostgresTodoRepository) RemoveTodo(id string) error {
	db, err := database.GetDatabaseConnectionFor(repository.ConnectionString)
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
