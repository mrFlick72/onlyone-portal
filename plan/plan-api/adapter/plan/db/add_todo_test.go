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

func TestAddANewTodo(t *testing.T) {
	aPlan := test.ANewPlan()
	planId, err := repo.CreateNewPlan(aPlan)
	require.NoError(t, err)

	aTodo := test.ANewTodoWith("A Content")
	require.NoError(t, repo.AddTodo(planId, aTodo))

	anotherTodo := test.ANewTodoWith("Another Content")
	require.NoError(t, repo.AddTodo(planId, anotherTodo))

	actual, err := repo.GetPlan(planId, aPlan.UserName)
	require.NoError(t, err)
	assert.Equal(t, plan.Plan{
		Id:       planId,
		UserName: aPlan.UserName,
		Title:    aPlan.Title,
		Date:     clock.ToDay(),
		Todos:    []*plan.Todo{&aTodo, &anotherTodo},
	}, *actual)
}
