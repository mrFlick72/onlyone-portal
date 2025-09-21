package plan

import (
	"fmt"
	"sort"
	"testing"
)

var (
	testableMySqlTodoRepository = MySqlTodoRepository{ConnectionString: testableDatabaseConnectionString}
)

func TestMySqlTodoRepository_SaveTodo(t *testing.T) {
	todo := aNewTodo()

	err := testableMySqlTodoRepository.SaveTodo(&todo)

	assertThatNoErrorFor(t, err, "some errors occurs during the insert query")
	clearDatabase()
}

func TestMySqlTodoRepository_GetTodo(t *testing.T) {
	expected := aNewTodo()
	err := testableMySqlTodoRepository.SaveTodo(&expected)
	assertThatNoErrorFor(t, err, "some errors occurs during the insert query")

	actual, err := testableMySqlTodoRepository.GetTodo(expected.Id)
	t.Log("actual ", actual)
	assertThatNoErrorFor(t, err, "some errors occurs during the find one query")

	if expected != *actual {
		t.Log("expected ", expected)
		t.Log("actual ", actual)
		t.Error("the retrieved todo is not wat we expect")
	}

	clearDatabase()
}

func TestMySqlTodoRepository_GetTodoNotFound(t *testing.T) {
	actual, err := testableMySqlTodoRepository.GetTodo("123456789987")
	fmt.Println("actual ", actual)
	assertThatNoErrorFor(t, err, "some errors occurs during the find one query")

	if actual != nil {
		fmt.Println("actual ", actual)
		t.Error("the retrieved todo is not wat we expect")
	}

	clearDatabase()
}

func TestMySqlTodoRepository_GetAllTodo(t *testing.T) {
	aTodo := aNewTodo()
	anotherTodo := aNewTodo()

	testableMySqlTodoRepository.SaveTodo(&aTodo)
	testableMySqlTodoRepository.SaveTodo(&anotherTodo)

	expected := orderedTodoListById(aTodo, anotherTodo)

	actual, err := testableMySqlTodoRepository.GetAllTodo("user-name")
	assertThatNoErrorFor(t, err, "some errors occurs during the select query")

	assertEqualityFor(t, expected[0], actual[0])
	assertEqualityFor(t, expected[1], actual[1])

	clearDatabase()
}

func TestMySqlTodoRepository_RemoveTodo(t *testing.T) {
	aTodo := aNewTodo()
	anotherTodo := aNewTodo()

	testableMySqlTodoRepository.SaveTodo(&aTodo)
	testableMySqlTodoRepository.SaveTodo(&anotherTodo)

	expected := []Todo{aTodo}

	sort.Slice(expected, func(p, q int) bool {
		return expected[p].Id < expected[q].Id
	})

	testableMySqlTodoRepository.RemoveTodo(anotherTodo.Id)
	actual, err := testableMySqlTodoRepository.GetAllTodo("user-name")
	assertThatNoErrorFor(t, err, "some errors occurs during the select query")

	assertEqualityFor(t, expected[0], actual[0])
	clearDatabase()
}
