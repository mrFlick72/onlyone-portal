package plan

import (
	"database/sql"
	"time"
)


func BuildTodos(rows *sql.Rows, result []*Todo) []*Todo {
	for rows.Next() {
		var id, content, username string
		var date time.Time
		err := rows.Scan(&id, &username, &date, &content)
		logger.LogErrorFor(err)

		result = append(result, &Todo{
			Id:       id,
			UserName: username,
			Date:     date,
			Content:  content,
		})
	}
	return result
}
