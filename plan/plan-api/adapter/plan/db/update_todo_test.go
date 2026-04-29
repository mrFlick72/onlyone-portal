//go:build test

package db

import (
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateTodo(t *testing.T) {
	aPlan := test.ANewPlan()
	planId, err := repo.CreateNewPlan(aPlan)
	require.NoError(t, err)

	aTodo := test.ANewTodoWith("Original Content")
	require.NoError(t, repo.AddTodo(planId, aTodo))

	updatedTodo := plan.Todo{Id: aTodo.Id, UserName: aTodo.UserName, Date: clock.ToDay(), Content: "Updated Content"}
	require.NoError(t, repo.UpdateTodo(planId, updatedTodo))

	actual, err := repo.GetPlan(planId, aPlan.UserName)
	require.NoError(t, err)
	assert.Len(t, actual.Todos, 1)
	assert.Equal(t, "Updated Content", actual.Todos[0].Content)
}
