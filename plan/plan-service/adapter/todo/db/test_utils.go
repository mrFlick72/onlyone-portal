//go:build test

package db

import (
	"database/sql"
	"fmt"
	"sort"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/todo"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
)

const testDSN = "host=localhost dbname=todo user=root password=root sslmode=disable"

func clearDatabase() {
	conn, _ := sql.Open("postgres", testDSN)
	defer conn.Close()
	cleanTable("todo", conn)
}

func cleanTable(table string, conn *sql.DB) {
	if _, err := conn.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)); err != nil {
		logger.LogErrorFor(err)
	}
}

func assertEqualityFor(t *testing.T, expected todo.Todo, actual *todo.Todo) {
	t.Helper()
	if expected != *actual {
		t.Errorf("expected %+v, got %+v", expected, *actual)
	}
}

func assertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: %v", msg, err)
	}
}

func aNewTodo() todo.Todo {
	id, _ := uuid.NewRandom()
	return todo.Todo{
		Id:       id.String(),
		Content:  "it is a todo",
		UserName: "user-name",
		Date:     clock.ToDay(),
	}
}

func orderedByID(a, b todo.Todo) []todo.Todo {
	list := []todo.Todo{a, b}
	sort.Slice(list, func(i, j int) bool { return list[i].Id < list[j].Id })
	return list
}
