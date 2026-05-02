//go:build test

package testutils

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

var JwkStr = `{
		"kty": "RSA",
		"kid": "test-key",
		"alg": "RS256",
		"n": "mmO0OvOPQ53HRxV4eHOkTTxLVfk6zcq8KAD86gbnydYBNO_Si4Q1twyvefd58-BaO4N4NCEA97QrYm57ThKCe8agLGwWPHhxgbu_SAuYQehXxkf4sWy7Q17kGFG5k5AfQGZBqTY-YaawQqLlF6ILVbWab_AoEF4yB7pI3AnNnXs",
		"e": "AQAB",
		"d": "RzsrI2vONJcuIyjPzVslehEQfRkhPWOFTjuudNc8yA25vs_LZ11XXx42M-KvXIqtdvngUsTLan2w6pgowcuecX3t_2wUx0GJJgARfkN7gsWIS3CyXZBEEMjLGVU4vHt5zNE3GJKo3hb1TwEiulpL_Ix6hfcTSJpEaBWrBxjxV-E",
		"p": "5EA0bi6ui1H1wsG85oc7i9O7UH58WPIK_ytzBWXFIwcaSFFBqqNYNnZaHFsMe4cbHSBgShWHO3UueGVgOKmB8Q",
		"q": "rSi7CosQZmj_RFIYW10ef7XTZsdpIdOXV9-1dThAJUvkslKiTfdU7T0IYYsJ2K58ekJqdpcoKAVLB2SZVvdqKw",
		"dp": "S9yjEHPng1qsShzGQgB0ZBbtTOWdQpq_2OuCAStACFJWA-8t2h8MNJ3FeWMxlOTkuBuIpVbeaX6bAV0ATBTaoQ",
		"dq": "ZssMJhkh1jm0d-FoVix0Y4oUAiqUzaDnciH6faiz47AnBnkporEV-HPH2ugII1qJyKZOvzHCg-eIf84HfWoI2w",
		"qi": "lyVz1HI2b1IjzOMENkmUTaVEO6DM6usZi3c3_MobUUM05yyBhnHtPjWzqWn1uJ_Gt5bkJDdcpfvmkPAhKWEU9Q"
	}`

func SetupTestWebServer(middleware []gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middleware...)
	RegisterEndpoints(router)
	return router
}

func RegisterEndpoints(r *gin.Engine) *gin.Engine {

	// GET /api/tags — return all tags as JSON
	r.GET("/api/sample", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "sample endpoint"})
	})

	return r
}

func CreateTestToken(authorities []string) string {
	return createTestTokenWithExpiry(authorities, time.Now().Add(time.Hour))
}

func CreateExpiredTestToken(authorities []string) string {
	return createTestTokenWithExpiry(authorities, time.Now().Add(-time.Minute))
}

// CreateTokenSignedByForeignKey signs a token with a freshly generated RSA
// key that is NOT in the test JWKS — used to prove the middleware verifies
// signatures.
func CreateTokenSignedByForeignKey(authorities []string) string {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("foreign key gen: %v", err)
	}
	jwkey, err := jwk.Import(priv)
	if err != nil {
		log.Fatalf("foreign key import: %v", err)
	}
	jwkey.Set(jwk.KeyIDKey, "test-key") // same kid as the trusted JWKS, but different key material
	jwkey.Set(jwk.AlgorithmKey, jwa.RS256())

	now := time.Now()
	t := jwt.New()
	t.Set(jwt.IssuedAtKey, now)
	t.Set(jwt.ExpirationKey, now.Add(time.Hour))
	t.Set("user_name", "test_user")
	t.Set("authorities", authorities)

	signed, err := jwt.Sign(t, jwt.WithKey(jwa.RS256(), jwkey))
	if err != nil {
		log.Fatalf("foreign key sign: %v", err)
	}
	return fmt.Sprintf("Bearer %s", string(signed))
}

func createTestTokenWithExpiry(authorities []string, exp time.Time) string {
	now := time.Now()
	t := jwt.New()
	t.Set(jwt.SubjectKey, `https://github.com/lestrrat-go/jwx/v3/jwt`)
	t.Set(jwt.AudienceKey, `Golang Users`)
	t.Set(jwt.IssuedAtKey, now)
	t.Set(jwt.ExpirationKey, exp)
	t.Set("user_name", "test_user")
	t.Set("authorities", authorities)
	jwkey, err := jwk.ParseKey([]byte(JwkStr))
	log.Println("error key parsing", err)

	signed, err := jwt.Sign(t, jwt.WithKey(jwa.RS256(), jwkey))

	log.Println("error token sign", err)
	log.Println("string(token)")
	log.Println(string(signed))
	return fmt.Sprintf("Bearer %s", string(signed))
}
