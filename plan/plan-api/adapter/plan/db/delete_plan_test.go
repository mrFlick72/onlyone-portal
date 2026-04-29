//go:build test

package db

import (
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
)

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
