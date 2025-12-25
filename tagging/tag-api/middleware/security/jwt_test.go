package security

import (
	"testing"

	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/stretchr/testify/assert"
)

func TestGetClaimListFromToken(t *testing.T) {
	tok, err := jwt.NewBuilder().Claim("authorities", []interface{}{"a","b"}).Build()
	assert.NoError(t, err)

	got := getClaimListFromToken(tok, "authorities")
	assert.NotNil(t, got)
	assert.Equal(t, []string{"a", "b"}, *got)
}

func TestGetClaimFromToken(t *testing.T) {
	tok, err := jwt.NewBuilder().Claim("claim_name", "a_claim").Build()
	assert.NoError(t, err)

	got := getClaimFromToken(tok, "claim_name")
	assert.NotNil(t, got)
	assert.Equal(t, "a_claim", *got)
}
