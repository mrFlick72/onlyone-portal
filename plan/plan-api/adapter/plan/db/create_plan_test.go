//go:build test

package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNewPlan(t *testing.T) {
	planId, err := repo.CreateNewPlan(test.ANewPlan())
	require.NoError(t, err)
	assert.NotEmpty(t, planId)
	_, err = uuid.Parse(planId)
	assert.NoError(t, err, "returned id should be a valid UUID")
}
