//go:build test

package db

import (
	"testing"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEmptyPlan(t *testing.T) {
	p := test.ANewPlan()
	planId, err := repo.CreateNewPlan(p)
	require.NoError(t, err)

	p.Id = planId
	actual, err := repo.GetPlan(planId, p.UserName)
	require.NoError(t, err)
	assert.Equal(t, p, *actual)
}

func TestGetPlanNotFound(t *testing.T) {
	_, err := repo.GetPlan("non-existent-id", "user-name")
	assert.Error(t, err)
}
