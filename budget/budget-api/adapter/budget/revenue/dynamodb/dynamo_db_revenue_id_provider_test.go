package dynamodb

import (
	"encoding/base64"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/budget/revenue"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestRevenueIdProviderProducesPythonCompatibleId(t *testing.T) {
	provider := &DynamoDbRevenueIdProvider{
		SaltGenerator: func() string { return "A_SALT" },
	}

	r := &revenue.Revenue{
		UserName: "USER",
		Date:     testutils.SafeDateFor("06/01/2018"),
	}

	id := provider.GenerateIdFor(r)

	expectedPk := base64.StdEncoding.EncodeToString([]byte("2018_USER"))
	expectedRk := base64.StdEncoding.EncodeToString([]byte("1_6_A_SALT"))
	assert.Equal(t, expectedPk+"-"+expectedRk, id)
}
