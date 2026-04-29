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

func TestGetAllPlanByWithNoPlans(t *testing.T) {
	plans, err := repo.GetAllPlanBy("unknown-user")
	require.NoError(t, err)
	assert.Empty(t, plans)
}

func TestGetAllPlanBy(t *testing.T) {
	userName := "all-plans-user"

	firstPlanId, err := repo.CreateNewPlan(plan.Plan{Title: "first plan", UserName: userName, Date: clock.ToDay(), Todos: []*plan.Todo{}})
	require.NoError(t, err)
	require.NoError(t, repo.AddTodo(firstPlanId, test.ANewTodoWith("A Content")))

	secondPlanId, err := repo.CreateNewPlan(plan.Plan{Title: "second plan", UserName: userName, Date: clock.ToDay(), Todos: []*plan.Todo{}})
	require.NoError(t, err)

	plans, err := repo.GetAllPlanBy(userName)
	require.NoError(t, err)
	assert.Len(t, plans, 2)
	assert.Equal(t, plan.Plan{Id: firstPlanId, Title: "first plan", UserName: userName, Date: clock.ToDay(), Todos: []*plan.Todo{}}, *plans[0])
	assert.Equal(t, plan.Plan{Id: secondPlanId, Title: "second plan", UserName: userName, Date: clock.ToDay(), Todos: []*plan.Todo{}}, *plans[1])
}
