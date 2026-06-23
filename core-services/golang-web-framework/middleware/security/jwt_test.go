package security

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGetClaimListFromToken(t *testing.T) {
	tok, err := jwt.NewBuilder().Claim("authorities", []interface{}{"a", "b"}).Build()
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

func TestWhenARequestIsAuthorized(t *testing.T) {
	router := SetupAuthenticatedTestWebServer([]string{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/sample", nil)
	req.Header.Set("Authorization", testutils.CreateTestToken([]string{"ROLE_USER", "ROLE_ADMIN"}))
	router.ServeHTTP(w, req)

	expected, _ := json.Marshal(map[string]string{"message": "sample endpoint"})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expected), w.Body.String())
}

func TestWhenARequestIsNoAuthorized(t *testing.T) {
	router := SetupAuthenticatedTestWebServer([]string{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/sample", nil)
	req.Header.Set("Authorization", testutils.CreateTestToken([]string{}))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// A token contains wrong authorities
func TestWhenTokenContainsWrongAuthorities(t *testing.T) {
	router := SetupAuthenticatedTestWebServer([]string{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/sample", nil)
	req.Header.Set("Authorization", testutils.CreateTestToken([]string{"WRONG"}))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// jwt.Parse validates `exp` by default; an expired token must be rejected
// without an explicit time check in the middleware.
func TestWhenATokenIsExpired(t *testing.T) {
	router := SetupAuthenticatedTestWebServer([]string{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/sample", nil)
	req.Header.Set("Authorization", testutils.CreateExpiredTestToken([]string{"ROLE_USER"}))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// A request without an Authorization header must be rejected — and must not
// crash the goroutine via the old `authorization[7:]` slice panic.
func TestWhenAuthorizationHeaderIsMissing(t *testing.T) {
	router := SetupAuthenticatedTestWebServer([]string{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/sample", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// A token signed by a key that is not in the JWKS must be rejected — proves
// signature verification is actually enabled (was previously a no-op via
// jwt.WithVerify(false)).
func TestWhenTokenIsSignedByForeignKey(t *testing.T) {
	router := SetupAuthenticatedTestWebServer([]string{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/sample", nil)
	req.Header.Set("Authorization", testutils.CreateTokenSignedByForeignKey([]string{"ROLE_USER"}))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestWhenAuthorizationEvaluationIsSkipped(t *testing.T) {
	router := SetupAuthenticatedTestWebServer([]string{"/api/sample"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/sample", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWhenAuthorizationEvaluationIsSkippedWithWildcard(t *testing.T) {
	router := SetupAuthenticatedTestWebServer([]string{"/api/*"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/sample", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func SetupAuthenticatedTestWebServer(ignored []string) *gin.Engine {
	keySet, _ := jwk.ParseString(testutils.JwkStr)
	middleware := NewOAuth2Middleware(keySet, "ROLE_USER", ignored)

	router := testutils.SetupTestWebServer([]gin.HandlerFunc{middleware})

	return router
}
