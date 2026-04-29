//go:build test

package test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
	"github.com/stretchr/testify/assert"
)

var logger = logging.GetLoggerInstanceForComponentByTypeName("TestUtils")

const TestDSN = "host=localhost dbname=postgres user=postgres password=postgres sslmode=disable"

func ClearDatabase() {
	conn, _ := sql.Open("postgres", TestDSN)
	defer conn.Close()
	InitDatabase(conn)
	if _, err := conn.Exec("TRUNCATE TABLE todo, plan"); err != nil {
		logger.LogErrorFor(err)
	}
}

func InitDatabase(conn *sql.DB) {
	content, err := os.ReadFile("../../../scripts/init.sql")
	if err != nil {
		logger.LogErrorfFor("File not found: %v", err)
	}
	if _, err := conn.Exec(string(content)); err != nil {
		logger.LogErrorfFor("Sql execution error: %v", err)
		panic(err)
	}
}

func ANewPlan() plan.Plan {
	return plan.Plan{
		Title:    "a test plan",
		UserName: "user-name",
		Date:     clock.ToDay(),
		Todos:    []*plan.Todo{},
	}
}

func ANewTodoWith(content string) plan.Todo {
	return plan.Todo{
		Id:       uuid.New().String(),
		UserName: "user-name",
		Date:     clock.ToDay(),
		Content:  content,
	}
}

func AssertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: %v", msg, err)
	}
}

func AssertValidUUID(t *testing.T, id string) {
	t.Helper()
	assert.NotEmpty(t, id)
	_, err := uuid.Parse(id)
	assert.NoError(t, err, "returned id should be a valid UUID")
}

func AssertEqualPlan(t *testing.T, expected plan.Plan, actual *plan.Plan) {
	t.Helper()
	assert.Equal(t, expected, *actual)
}
