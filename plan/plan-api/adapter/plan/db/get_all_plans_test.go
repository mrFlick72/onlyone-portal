//go:build test

package db

import (
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPlanByWithNoPlans(t *testing.T) {
	plans, err := repo.GetAllPlanBy("unknown-user")
	test.AssertNoError(t, err, "get all failed")
	assert.Equal(t, make([]*plan.Plan, 0), plans)
}

func TestGetAllPlanBy(t *testing.T) {
	userName := "all-plans-user"

	firstPlanId, err := repo.CreateNewPlan(plan.Plan{Title: "first plan", UserName: userName, Date: clock.ToDay(), Todos: []*plan.Todo{}})
	test.AssertNoError(t, err, "insert first plan failed")

	aTodo := test.ANewTodoWith("A Content")
	test.AssertNoError(t, repo.AddTodo(firstPlanId, aTodo), "adding todo to first plan failed")

	secondPlanId, err := repo.CreateNewPlan(plan.Plan{Title: "second plan", UserName: userName, Date: clock.ToDay(), Todos: []*plan.Todo{}})
	test.AssertNoError(t, err, "insert second plan failed")

	plans, err := repo.GetAllPlanBy(userName)
	test.AssertNoError(t, err, "get all failed")
	assert.Equal(t, 2, len(plans))

	test.AssertEqualPlan(t, plan.Plan{Id: firstPlanId, Title: "first plan", UserName: userName, Date: clock.ToDay(), Todos: []*plan.Todo{}}, plans[0])
	test.AssertEqualPlan(t, plan.Plan{Id: secondPlanId, Title: "second plan", UserName: userName, Date: clock.ToDay(), Todos: []*plan.Todo{}}, plans[1])
}
