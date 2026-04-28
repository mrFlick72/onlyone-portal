//go:build test

package db

import (
	"fmt"
	domainplan "github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/todo"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
	"os"
	"testing"
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
