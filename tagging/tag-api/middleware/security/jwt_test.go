package security

import (
    "encoding/json"
    "testing"

    "github.com/lestrrat-go/jwx/v3/jwt"
    "github.com/stretchr/testify/assert"
)

func TestGetClaimListFromToken_WithStringSlice(t *testing.T) {
    tok, err := jwt.NewBuilder().Claim("authorities", []string{"a", "b"}).Build()
    assert.NoError(t, err)

    got := getClaimListFromToken(tok, "authorities")
    assert.NotNil(t, got)
    assert.Equal(t, []string{"a", "b"}, *got)
}

func TestGetClaimListFromToken_WithInterfaceSlice(t *testing.T) {
    // []interface{} with mixed types
    val := []interface{}{"x", "y", 123}
    tok, err := jwt.NewBuilder().Claim("authorities", val).Build()
    assert.NoError(t, err)

    got := getClaimListFromToken(tok, "authorities")
    assert.NotNil(t, got)
    assert.Equal(t, []string{"x", "y", "123"}, *got)
}

func TestGetClaimListFromToken_WithRawJSON(t *testing.T) {
    // Some token issuers might embed the array as a raw json string
    rawArr, _ := json.Marshal([]string{"r1", "r2"})
    tok, err := jwt.NewBuilder().Claim("authorities", json.RawMessage(rawArr)).Build()
    assert.NoError(t, err)

    got := getClaimListFromToken(tok, "authorities")
    assert.NotNil(t, got)
    assert.Equal(t, []string{"r1", "r2"}, *got)
}
