//go:build test

package db

import (
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTodoIsCreatedInTodoStatus(t *testing.T) {
	aPlan := test.ANewPlan()
	planId, err := repo.CreateNewPlan(aPlan)
	require.NoError(t, err)

	aTodo := test.ANewTodoWith("status-default")
	require.NoError(t, repo.AddTodo(planId, aTodo))

	actual, err := repo.GetPlan(planId, aPlan.UserName)
	require.NoError(t, err)
	require.Len(t, actual.Todos, 1)
	assert.Equal(t, plan.StatusTodo, actual.Todos[0].Status)
}

func TestUpdateTodoStatusPersistsTransition(t *testing.T) {
	aPlan := test.ANewPlan()
	planId, err := repo.CreateNewPlan(aPlan)
	require.NoError(t, err)

	aTodo := test.ANewTodoWith("transition-content")
	require.NoError(t, repo.AddTodo(planId, aTodo))

	require.NoError(t, repo.UpdateTodoStatus(planId, aTodo.Id, plan.StatusInProgress))

	actual, err := repo.GetPlan(planId, aPlan.UserName)
	require.NoError(t, err)
	require.Len(t, actual.Todos, 1)
	assert.Equal(t, plan.StatusInProgress, actual.Todos[0].Status)
}

func TestUpdateTodoLeavesStatusUnchanged(t *testing.T) {
	aPlan := test.ANewPlan()
	planId, err := repo.CreateNewPlan(aPlan)
	require.NoError(t, err)

	aTodo := test.ANewTodoWith("before-update")
	require.NoError(t, repo.AddTodo(planId, aTodo))
	require.NoError(t, repo.UpdateTodoStatus(planId, aTodo.Id, plan.StatusInProgress))

	updated := plan.Todo{Id: aTodo.Id, UserName: aTodo.UserName, Date: aTodo.Date, Content: "after-update"}
	require.NoError(t, repo.UpdateTodo(planId, updated))

	actual, err := repo.GetPlan(planId, aPlan.UserName)
	require.NoError(t, err)
	require.Len(t, actual.Todos, 1)
	assert.Equal(t, "after-update", actual.Todos[0].Content)
	assert.Equal(t, plan.StatusInProgress, actual.Todos[0].Status, "UpdateTodo must not reset status")
}
