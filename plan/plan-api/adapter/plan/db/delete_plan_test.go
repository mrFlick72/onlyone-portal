//go:build test

package db

import (
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeletePlan(t *testing.T) {
	aPlan := test.ANewPlan()
	planId, err := repo.CreateNewPlan(aPlan)
	require.NoError(t, err)

	require.NoError(t, repo.AddTodo(planId, test.ANewTodoWith("A Content")))
	require.NoError(t, repo.DeletePlan(planId, aPlan.UserName))

	_, err = repo.GetPlan(planId, aPlan.UserName)
	assert.Error(t, err)
}
