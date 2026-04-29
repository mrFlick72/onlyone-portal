//go:build test

package db

import (
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
)

func TestCreateNewPlan(t *testing.T) {
	plan := test.ANewPlan()
	planId, err := repo.CreateNewPlan(plan)
	test.AssertNoError(t, err, "insert failed")
	test.AssertValidUUID(t, planId)
}
