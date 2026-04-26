//go:build test

package plan

import (
	"database/sql"
	"fmt"
	"sort"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/pkg/clock"
)

var testableDatabaseConnectionString = "host=localhost dbname=todo user=root password=root sslmode=disable"

func clearDatabase() {
	connection, _ := sql.Open("postgres", testableDatabaseConnectionString)
	cleanTable("todo", connection)
	cleanTable("plan", connection)
}

func cleanTable(table string, connection *sql.DB) {
	_, err := connection.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
	if err != nil {
		logger.LogErrorFor(err)
	}
}

func assertEqualityFor(t *testing.T, expected Todo, actual *Todo) {
	if expected != *actual {
		t.Error("expected: ", expected)
		t.Error("actual: ", actual)
		t.Error("the retrieved todo is not what we expect")
	}
}

func assertThatNoErrorFor(t *testing.T, err error, errorMessage string) {
	if err != nil {
		t.Log(errorMessage)
	}
}

func aNewTodo() Todo {
	random, _ := uuid.NewRandom()
	return Todo{
		Id:       random.String(),
		Content:  "it is a todo",
		UserName: "user-name",
		Date:     clock.ToDay(),
	}
}

func orderedTodoListById(aTodo Todo, anotherTodo Todo) []Todo {
	expected := []Todo{aTodo, anotherTodo}
	sort.Slice(expected, func(p, q int) bool {
		return expected[p].Id < expected[q].Id
	})
	return expected
}
