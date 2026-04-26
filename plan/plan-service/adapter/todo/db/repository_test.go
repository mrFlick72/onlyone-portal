//go:build test

package db

import (
	"sort"
	"testing"
)

var repo = TodoPostgresRepository{ConnectionString: testDSN}

func TestSaveTodo(t *testing.T) {
	clearDatabase()
	todo := aNewTodo()

	err := repo.SaveTodo(&todo)

	assertNoError(t, err, "insert failed")
	clearDatabase()
}

func TestGetTodo(t *testing.T) {
	clearDatabase()
	expected := aNewTodo()
	assertNoError(t, repo.SaveTodo(&expected), "insert failed")

	actual, err := repo.GetTodo(expected.Id)
	assertNoError(t, err, "get failed")
	assertEqualityFor(t, expected, actual)

	clearDatabase()
}

func TestGetTodoNotFound(t *testing.T) {
	clearDatabase()
	actual, _ := repo.GetTodo("non-existent-id")
	if actual != nil {
		t.Error("expected nil for missing todo")
	}
}

func TestGetAllTodo(t *testing.T) {
	clearDatabase()
	a, b := aNewTodo(), aNewTodo()
	repo.SaveTodo(&a)
	repo.SaveTodo(&b)
	expected := orderedByID(a, b)

	actual, err := repo.GetAllTodo("user-name")
	assertNoError(t, err, "get all failed")

	sort.Slice(actual, func(i, j int) bool { return actual[i].Id < actual[j].Id })
	assertEqualityFor(t, expected[0], actual[0])
	assertEqualityFor(t, expected[1], actual[1])

	clearDatabase()
}

func TestRemoveTodo(t *testing.T) {
	clearDatabase()
	a, b := aNewTodo(), aNewTodo()
	repo.SaveTodo(&a)
	repo.SaveTodo(&b)

	repo.RemoveTodo(b.Id)
	actual, err := repo.GetAllTodo("user-name")
	assertNoError(t, err, "get after remove failed")

	if len(actual) != 1 {
		t.Fatalf("expected 1 todo after removal, got %d", len(actual))
	}
	assertEqualityFor(t, a, actual[0])
	clearDatabase()
}
