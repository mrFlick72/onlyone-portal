//go:build test

package db

import (
	"testing"
)

var repo = PlanPostgresRepository{ConnectionString: testDSN}

func TestCreateNewPlan(t *testing.T) {
	clearDatabase()
	plan := aNewPlan()
	planId, err := repo.CreateNewPlan(plan)
	assertNoError(t, err, "insert failed")
	assertValidUUID(t, planId)
	clearDatabase()
}
