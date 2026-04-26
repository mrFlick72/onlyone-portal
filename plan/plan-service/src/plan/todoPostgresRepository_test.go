//go:build test

package plan

import (
	"fmt"
	"sort"
	"testing"
)

var testablePostgresTodoRepository = PostgresTodoRepository{ConnectionString: testableDatabaseConnectionString}

func TestPostgresTodoRepository_SaveTodo(t *testing.T) {
	todo := aNewTodo()

	err := testablePostgresTodoRepository.SaveTodo(&todo)

	assertThatNoErrorFor(t, err, "some errors occurred during the insert query")
	clearDatabase()
}

func TestPostgresTodoRepository_GetTodo(t *testing.T) {
	expected := aNewTodo()
	err := testablePostgresTodoRepository.SaveTodo(&expected)
	assertThatNoErrorFor(t, err, "some errors occurred during the insert query")

	actual, err := testablePostgresTodoRepository.GetTodo(expected.Id)
	assertThatNoErrorFor(t, err, "some errors occurred during the find one query")

	if expected != *actual {
		t.Log("expected ", expected)
		t.Log("actual ", actual)
		t.Error("the retrieved todo is not what we expect")
	}

	clearDatabase()
}

func TestPostgresTodoRepository_GetTodoNotFound(t *testing.T) {
	actual, _ := testablePostgresTodoRepository.GetTodo("non-existent-id")
	fmt.Println("actual ", actual)

	if actual != nil {
		t.Error("expected nil for missing todo")
	}
}

func TestPostgresTodoRepository_GetAllTodo(t *testing.T) {
	aTodo := aNewTodo()
	anotherTodo := aNewTodo()

	testablePostgresTodoRepository.SaveTodo(&aTodo)
	testablePostgresTodoRepository.SaveTodo(&anotherTodo)

	expected := orderedTodoListById(aTodo, anotherTodo)

	actual, err := testablePostgresTodoRepository.GetAllTodo("user-name")
	assertThatNoErrorFor(t, err, "some errors occurred during the select query")

	sort.Slice(actual, func(p, q int) bool {
		return actual[p].Id < actual[q].Id
	})

	assertEqualityFor(t, expected[0], actual[0])
	assertEqualityFor(t, expected[1], actual[1])

	clearDatabase()
}

func TestPostgresTodoRepository_RemoveTodo(t *testing.T) {
	aTodo := aNewTodo()
	anotherTodo := aNewTodo()

	testablePostgresTodoRepository.SaveTodo(&aTodo)
	testablePostgresTodoRepository.SaveTodo(&anotherTodo)

	testablePostgresTodoRepository.RemoveTodo(anotherTodo.Id)
	actual, err := testablePostgresTodoRepository.GetAllTodo("user-name")
	assertThatNoErrorFor(t, err, "some errors occurred during the select query")

	if len(actual) != 1 {
		t.Errorf("expected 1 todo after removal, got %d", len(actual))
	}
	assertEqualityFor(t, aTodo, actual[0])
	clearDatabase()
}
