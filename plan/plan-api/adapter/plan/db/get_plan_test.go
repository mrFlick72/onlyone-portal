//go:build test

package db

import (
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
)

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
