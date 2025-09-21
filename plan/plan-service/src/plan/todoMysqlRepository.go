package plan

import (
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/pkg/database"
)

type MySqlTodoRepository struct {
	ConnectionString string
}

func (repository *MySqlTodoRepository) GetAllTodo(userName string) ([]*Todo, error) {
	result := make([]*Todo, 0)

	db, err := database.GetDatabaseConnectionFor(repository.ConnectionString)

	query, _ := db.Prepare("SELECT id, user_name as username, date, content FROM TODO where user_name = ?")
	rows, err := query.Query(userName)
	logger.LogErrorFor(err)
	result = BuildTodos(rows, result)

	database.CloseResources(rows, query, db)
	return result, err
}

func (repository *MySqlTodoRepository) GetTodo(id string) (*Todo, error) {
	var selectAll []*Todo

	db, err := database.GetDatabaseConnectionFor(repository.ConnectionString)

	query, _ := db.Prepare("SELECT id, user_name as username, date, content FROM TODO WHERE ID=?")
	rows, err := query.Query(id)
	logger.LogErrorFor(err)

	selectAll = BuildTodos(rows, selectAll)

	database.CloseResources(rows, query, db)

	if len(selectAll) > 0 {
		return selectAll[0], err
	} else {
		err := errors.New("todo with id " + id + " not found")
		logger.LogErrorFor(err)
		return nil, err
	}

}

func (repository *MySqlTodoRepository) SaveTodo(todo *Todo) error {
	db, err := database.GetDatabaseConnectionFor(repository.ConnectionString)

	query, _ := db.Prepare("INSERT into TODO (id, user_name, date, content) VALUES (?, ?, ?, ?)")
	rows, err := query.Query(todo.Id, todo.UserName, todo.Date, todo.Content)
	if err == nil {
		logger.LogErrorFor(err)
	}

	database.CloseResources(rows, query, db)

	return err
}

func (repository *MySqlTodoRepository) RemoveTodo(id string) error {
	db, err := database.GetDatabaseConnectionFor(repository.ConnectionString)

	query, _ := db.Prepare("DELETE FROM TODO WHERE id=?")
	rows, err := query.Query(id)
	logger.LogErrorFor(err)

	database.CloseResources(rows, query, db)

	return err
}
