package db

import (
	"database/sql"
	"time"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/todo"
)

func buildTodos(rows *sql.Rows, result []*todo.Todo) []*todo.Todo {
	for rows.Next() {
		var id, content, username string
		var date time.Time
		if err := rows.Scan(&id, &username, &date, &content); err != nil {
			logger.LogErrorFor(err)
			continue
		}
		result = append(result, &todo.Todo{
			Id:       id,
			UserName: username,
			Date:     date,
			Content:  content,
		})
	}
	return result
}
