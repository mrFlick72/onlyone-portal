package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
)

var logger = logging.GetLoggerInstance()

func GetDatabaseConnectionFor(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	logger.LogErrorFor(err)
	return db, err
}

func CloseResources(rows *sql.Rows, query *sql.Stmt, db *sql.DB) {
	if rows != nil {
		defer rows.Close()
	}
	if query != nil {
		defer query.Close()
	}
	if db != nil {
		defer db.Close()
	}
}
