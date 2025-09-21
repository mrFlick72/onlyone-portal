package database

import (
	"database/sql"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/pkg/logging"
)

var logger = logging.GetLoggerInstance()

func GetDatabaseConnectionFor(connectionStrung string) (*sql.DB, error) {
	database, err := sql.Open("mysql", connectionStrung)
	logger.LogErrorFor(err)
	return database, err
}

func CloseResources(rows *sql.Rows, query *sql.Stmt, database *sql.DB) {
	if rows != nil {
		defer rows.Close()
	}

	if query != nil {
		defer query.Close()
	}

	if database != nil {
		defer database.Close()
	}
}
