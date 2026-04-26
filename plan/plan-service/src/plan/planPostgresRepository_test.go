//go:build test

package plan

import "testing"

var testablePostgresPlanRepository = PostgresPlanRepository{ConnectionString: testableDatabaseConnectionString}

func TestPostgresPlanRepository_CreateNewPlan(t *testing.T) {
	actual, err := testablePostgresPlanRepository.CreateNewPlan(Plan{})
	if err != nil {
		t.Fail()
	}

	_ = actual
}
