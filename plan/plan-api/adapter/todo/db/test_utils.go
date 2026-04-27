//go:build test

package db

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"sort"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/todo"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
)

const testDSN = "host=localhost dbname=postgres user=postgres password=postgres sslmode=disable"

func clearDatabase() {
	conn, _ := sql.Open("postgres", testDSN)
	defer conn.Close()
	initDatabase(conn)
	cleanTable("todo", conn)
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

func cleanTable(table string, conn *sql.DB) {
	if _, err := conn.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)); err != nil {
		logger.LogErrorFor(err)
	}
}

func assertEqualityFor(t *testing.T, expected todo.Todo, actual *todo.Todo) {
	t.Helper()
	assert.Equal(t, expected, *actual)
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
