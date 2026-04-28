//go:build test

package db

import (
	"database/sql"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/todo"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	domainplan "github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
	"github.com/stretchr/testify/assert"
)

const testDSN = "host=localhost dbname=postgres user=postgres password=postgres sslmode=disable"

func clearDatabase() {
	conn, _ := sql.Open("postgres", testDSN)
	defer conn.Close()
	initDatabase(conn)
	if _, err := conn.Exec("TRUNCATE TABLE todo, plan"); err != nil {
		logger.LogErrorFor(err)
	}
}

func initDatabase(conn *sql.DB) {
	content, err := os.ReadFile("../../../scripts/init.sql")
	if err != nil {
		logger.LogErrorfFor("File not found: %v", err)
	}
	if _, err := conn.Exec(string(content)); err != nil {
		logger.LogErrorfFor("Sql execution error: %v", err)
		panic(err)
	}
}

func aNewPlan() domainplan.Plan {
	return domainplan.Plan{
		Title:    "a test plan",
		UserName: "user-name",
		Date:     clock.ToDay(),
		Todos:    []*todo.Todo{},
	}
}

func aNewTodoWith(content string) todo.Todo {
	return todo.Todo{
		Id:       uuid.New().String(),
		UserName: "user-name",
		Date:     clock.ToDay(),
		Content:  content,
	}
}

func assertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: %v", msg, err)
	}
}

func assertValidUUID(t *testing.T, id string) {
	t.Helper()
	assert.NotEmpty(t, id)
	_, err := uuid.Parse(id)
	assert.NoError(t, err, "returned id should be a valid UUID")
}

func assertEqualPlan(t *testing.T, expected domainplan.Plan, actual *domainplan.Plan) {
	t.Helper()
	assert.Equal(t, expected, *actual)
}
