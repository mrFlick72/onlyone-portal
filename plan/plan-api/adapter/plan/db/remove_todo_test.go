//go:build test

package db

import (
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
)

func TestRemoveATodo(t *testing.T) {
	aPlan := test.ANewPlan()

	planId, err := repo.CreateNewPlan(aPlan)
	test.AssertNoError(t, err, "insert failed")

	aTodo := test.ANewTodoWith("A Content")
	test.AssertNoError(t, repo.AddTodo(planId, aTodo), "adding the todo failed")

	anotherTodo := test.ANewTodoWith("Another Content")
	test.AssertNoError(t, repo.AddTodo(planId, anotherTodo), "adding the todo failed")

	test.AssertNoError(t, repo.RemoveTodo(planId, anotherTodo.Id), "removing the todo failed")

	expected := plan.Plan{
		Id:       planId,
		UserName: aPlan.UserName,
		Title:    aPlan.Title,
		Date:     clock.ToDay(),
		Todos:    []*plan.Todo{&aTodo},
	}

	actual, err := repo.GetPlan(planId, aPlan.UserName)
	test.AssertNoError(t, err, "get failed")
	test.AssertEqualPlan(t, expected, actual)
}
