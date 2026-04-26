package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
)

var logger = logging.GetLoggerInstance()

func GetDatabaseConnectionFor(connectionString string) (*sql.DB, error) {
	database, err := sql.Open("postgres", connectionString)
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
