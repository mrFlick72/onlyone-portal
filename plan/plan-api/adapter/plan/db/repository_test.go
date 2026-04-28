//go:build test

package db

import (
	"fmt"
	"os"
	"testing"

	domainplan "github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/todo"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
	"github.com/stretchr/testify/assert"
)

var repo = PlanPostgresRepository{ConnectionString: testDSN}

func TestMain(m *testing.M) {
	// global setup
	fmt.Println("setup before all tests")
	clearDatabase()

	code := m.Run()

	// global teardown
	fmt.Println("cleanup after all tests")

	os.Exit(code)
}

func TestCreateNewPlan(t *testing.T) {
	plan := aNewPlan()
	planId, err := repo.CreateNewPlan(plan)
	assertNoError(t, err, "insert failed")
	assertValidUUID(t, planId)
}

func TestGetEmptyPlan(t *testing.T) {
	plan := aNewPlan()
	planId, err := repo.CreateNewPlan(plan)
	assertNoError(t, err, "insert failed")

	plan.Id = planId
	actual, err := repo.GetPlan(planId, plan.UserName)
	assertNoError(t, err, "get failed")
	assertEqualPlan(t, plan, actual)
}

func TestGetPlanNotFound(t *testing.T) {
	_, err := repo.GetPlan("non-existent-id", "user-name")
	if err == nil {
		t.Error("expected error for non-existent plan, got nil")
	}
}

func TestAddANewTodo(t *testing.T) {
	plan := aNewPlan()

	planId, err := repo.CreateNewPlan(plan)
	assertNoError(t, err, "insert failed")

	aTodo := aNewTodoWith("A Content")
	err = repo.AddTodo(planId, aTodo)
	assertNoError(t, err, "adding the todo failed")

	anotherTodo := aNewTodoWith("Another Content")
	err = repo.AddTodo(planId, anotherTodo)
	assertNoError(t, err, "adding the todo failed")

	expected := domainplan.Plan{
		Id:       planId,
		UserName: plan.UserName,
		Title:    plan.Title,
		Date:     clock.ToDay(),
		Todos:    []*todo.Todo{&aTodo, &anotherTodo},
	}

	actual, err := repo.GetPlan(planId, plan.UserName)
	assertNoError(t, err, "get failed")
	assertEqualPlan(t, expected, actual)
}

func TestGetAllPlanByWithNoPlans(t *testing.T) {
	plans, err := repo.GetAllPlanBy("unknown-user")
	assertNoError(t, err, "get all failed")
	assert.Equal(t, make([]*domainplan.Plan, 0), plans)
}

func TestGetAllPlanBy(t *testing.T) {
	userName := "all-plans-user"

	firstPlanId, err := repo.CreateNewPlan(domainplan.Plan{Title: "first plan", UserName: userName, Date: clock.ToDay(), Todos: []*todo.Todo{}})
	assertNoError(t, err, "insert first plan failed")

	aTodo := aNewTodoWith("A Content")
	assertNoError(t, repo.AddTodo(firstPlanId, aTodo), "adding todo to first plan failed")

	secondPlanId, err := repo.CreateNewPlan(domainplan.Plan{Title: "second plan", UserName: userName, Date: clock.ToDay(), Todos: []*todo.Todo{}})
	assertNoError(t, err, "insert second plan failed")

	plans, err := repo.GetAllPlanBy(userName)
	assertNoError(t, err, "get all failed")
	assert.Equal(t, 2, len(plans))

	assertEqualPlan(t, domainplan.Plan{Id: firstPlanId, Title: "first plan", UserName: userName, Date: clock.ToDay(), Todos: []*todo.Todo{}}, plans[0])
	assertEqualPlan(t, domainplan.Plan{Id: secondPlanId, Title: "second plan", UserName: userName, Date: clock.ToDay(), Todos: []*todo.Todo{}}, plans[1])
}

func TestRemoveANewTodo(t *testing.T) {
	plan := aNewPlan()

	planId, err := repo.CreateNewPlan(plan)
	assertNoError(t, err, "insert failed")

	aTodo := aNewTodoWith("A Content")
	err = repo.AddTodo(planId, aTodo)
	assertNoError(t, err, "adding the todo failed")

	anotherTodo := aNewTodoWith("Another Content")
	err = repo.AddTodo(planId, anotherTodo)
	assertNoError(t, err, "adding the todo failed")

	err = repo.RemoveTodo(planId, anotherTodo.Id)
	assertNoError(t, err, "removing the todo failed")

	expected := domainplan.Plan{
		Id:       planId,
		UserName: plan.UserName,
		Title:    plan.Title,
		Date:     clock.ToDay(),
		Todos:    []*todo.Todo{&aTodo},
	}

	actual, err := repo.GetPlan(planId, plan.UserName)
	assertNoError(t, err, "get failed")
	assertEqualPlan(t, expected, actual)
}
