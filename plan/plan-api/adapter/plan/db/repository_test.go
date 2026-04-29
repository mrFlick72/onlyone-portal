//go:build test

package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
)

var repo = PlanPostgresRepository{ConnectionString: test.TestDSN}

func TestMain(m *testing.M) {
	fmt.Println("setup before all tests")
	test.ClearDatabase()

	code := m.Run()

	fmt.Println("cleanup after all tests")
	os.Exit(code)
}
