//go:build test

package db

import (
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateTodo(t *testing.T) {
	aPlan := test.ANewPlan()
	planId, err := repo.CreateNewPlan(aPlan)
	test.AssertNoError(t, err, "insert plan failed")

	aTodo := test.ANewTodoWith("Original Content")
	test.AssertNoError(t, repo.AddTodo(planId, aTodo), "add todo failed")

	updatedTodo := plan.Todo{Id: aTodo.Id, UserName: aTodo.UserName, Date: clock.ToDay(), Content: "Updated Content"}
	test.AssertNoError(t, repo.UpdateTodo(planId, updatedTodo), "update todo failed")

	actual, err := repo.GetPlan(planId, aPlan.UserName)
	test.AssertNoError(t, err, "get plan failed")
	assert.Equal(t, 1, len(actual.Todos))
	assert.Equal(t, "Updated Content", actual.Todos[0].Content)
}
