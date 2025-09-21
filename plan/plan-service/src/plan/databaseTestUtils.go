package plan

import (
	"database/sql"
	"fmt"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/pkg/clock"
)

var testableDatabaseConnectionString = "root:root@tcp(localhost)/todo?parseTime=true"

func clearDatabase() {
	connection, _ := sql.Open("mysql", testableDatabaseConnectionString)
	cleanTable("TODO", connection)
	cleanTable("PLAN", connection)
}

func cleanTable(table string, connection *sql.DB) {
	_, error := connection.Exec(fmt.Sprintf("truncate table %s", table))
	if error != nil {
		logger.LogErrorFor(error)
	}
}

func assertEqualityFor(t *testing.T, expected Todo, actual *Todo) {
	if expected != *actual {
		t.Error("expected: ", expected)
		t.Error("actual: ", actual)
		t.Error("the retrieved todo is not wat we expect")
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
