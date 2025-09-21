package plan

import "testing"

var (
	testableMySqlPlanRepository = MySqlPlanRepository{ConnectionString: testableDatabaseConnectionString}
)

func TestMySqlPlanRepository_CreateNewPlan(t *testing.T) {

	actual, err := testableMySqlPlanRepository.CreateNewPlan(Plan{})
	if err != nil {
		t.Fail()
	}


	if &actual == nil {
		t.Fail()
	}
}
