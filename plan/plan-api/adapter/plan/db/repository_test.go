//go:build test

package db

import (
	"fmt"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
	"os"
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
	"github.com/stretchr/testify/assert"
)

var repo = PlanPostgresRepository{ConnectionString: test.TestDSN}

func TestMain(m *testing.M) {
	// global setup
	fmt.Println("setup before all tests")
	test.ClearDatabase()

	code := m.Run()

	// global teardown
	fmt.Println("cleanup after all tests")

	os.Exit(code)
}

func TestCreateNewPlan(t *testing.T) {
	plan := test.ANewPlan()
	planId, err := repo.CreateNewPlan(plan)
	test.AssertNoError(t, err, "insert failed")
	test.AssertValidUUID(t, planId)
}

func TestGetEmptyPlan(t *testing.T) {
	plan := test.ANewPlan()
	planId, err := repo.CreateNewPlan(plan)
	test.AssertNoError(t, err, "insert failed")

	plan.Id = planId
	actual, err := repo.GetPlan(planId, plan.UserName)
	test.AssertNoError(t, err, "get failed")
	test.AssertEqualPlan(t, plan, actual)
}

func TestGetPlanNotFound(t *testing.T) {
	_, err := repo.GetPlan("non-existent-id", "user-name")
	if err == nil {
		t.Error("expected error for non-existent plan, got nil")
	}
}

func TestAddANewTodo(t *testing.T) {
	aPlan := test.ANewPlan()

	planId, err := repo.CreateNewPlan(aPlan)
	test.AssertNoError(t, err, "insert failed")

	aTodo := test.ANewTodoWith("A Content")
	err = repo.AddTodo(planId, aTodo)
	test.AssertNoError(t, err, "adding the todo failed")

	anotherTodo := test.ANewTodoWith("Another Content")
	err = repo.AddTodo(planId, anotherTodo)
	test.AssertNoError(t, err, "adding the todo failed")

	expected := plan.Plan{
		Id:       planId,
		UserName: aPlan.UserName,
		Title:    aPlan.Title,
		Date:     clock.ToDay(),
		Todos:    []*plan.Todo{&aTodo, &anotherTodo},
	}

	actual, err := repo.GetPlan(planId, aPlan.UserName)
	test.AssertNoError(t, err, "get failed")
	test.AssertEqualPlan(t, expected, actual)
}

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

func TestDeletePlan(t *testing.T) {
	aPlan := test.ANewPlan()
	planId, err := repo.CreateNewPlan(aPlan)
	test.AssertNoError(t, err, "insert plan failed")

	aTodo := test.ANewTodoWith("A Content")
	test.AssertNoError(t, repo.AddTodo(planId, aTodo), "add todo failed")
	test.AssertNoError(t, repo.DeletePlan(planId, aPlan.UserName), "delete plan failed")

	_, err = repo.GetPlan(planId, aPlan.UserName)
	if err == nil {
		t.Error("expected error for deleted plan, got nil")
	}
}

func TestRemoveANewTodo(t *testing.T) {
	aPlan := test.ANewPlan()

	planId, err := repo.CreateNewPlan(aPlan)
	test.AssertNoError(t, err, "insert failed")

	aTodo := test.ANewTodoWith("A Content")
	err = repo.AddTodo(planId, aTodo)
	test.AssertNoError(t, err, "adding the todo failed")

	anotherTodo := test.ANewTodoWith("Another Content")
	err = repo.AddTodo(planId, anotherTodo)
	test.AssertNoError(t, err, "adding the todo failed")

	err = repo.RemoveTodo(planId, anotherTodo.Id)
	test.AssertNoError(t, err, "removing the todo failed")

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
