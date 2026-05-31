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
			"alg": "RS256",
			"d": "DS-VAKWRHFx2UzpGJgYwCvG6QdMmBOsaJG0MH7Te2xkkjAr5wq3C20e2RuxVgvrluUk8neFFXAYyQHqbsbvT4bEgX9VoDa8ByUFTsLgIb8Q3w-9gIuUvr5930fnZSkQAXWpgemp3vQh09LdAgWs0tqOsafK6T5LeRroQaenLOebzp0PIMcF8qTromO-4Gpa1Kw8kPBt8UFlvuXwRHh03sLO-NvplLO2Kk43pLg_15_1HS4jSixZS28rzVConMU-KdvVBAlbA4wUxS3g2ylQQzd9YlWW_bJTs4HAs-5FBKGd714QWv9AH4uGtnt29QyBKDts3Ltqf4yTKl8jdxD2z4Q",
			"dp": "YiyV0Wzut95PEm1kDcAYtqDEYHcN5MVDz-tf5mrvFbRS1NDmVVe_5-r8l_ClQTnCLAX56mKSDVBg95mhUbJLjX5_d3NW36K-bDKsQ0AmpuURSFoaMqvRuJTGifbCoE0pxpast1loOlDClhiPy9wUCvm27ljw3BOy0gW856yHcRE",
			"dq": "wG9njVaAM1yDKCR6DZIoG79ndQWuboukKW2vIqFjegaUVeolxE5vWIr1mqIF3t0_-CT0yB1bx4BTVZOn-bcHUP1OxAhkLRTWd5Kr2gsDv978vRcJnCPdyIUX9AZbKEo2LmcETCUyNB5djY1HRQqhd8ghyKJtZJgi_L54INBYo1E",
			"e": "AQAB",
			"kid": "test-key",
			"kty": "RSA",
			"n": "rxmYcy9OstuH5p0qOeFfiqd6dXwOXzeE1P7FSAldt7d3BDCX911OXgwJhoWriwzp5jEXcGnHZbXdGXLkVMkWNw_w5zublLt9JCXCgihj5QI-vy48APNPAPRkCCBvAu7TAoqR4OZsGcpWjpSL1EHVBCgmx8qV5ia255laZfb0oyb_kWQ4DDJq6Qi1eIuhMVubrheCqi2bOqM_pOt6_LuOksSzWutp5XO0pg9vf4VzbMripqurwqxetK0Zx-1VgCxfp1USTWguwB0xK-6a1EjdGua-RVVQah2BZpV9g2O1PfWaeK88FX6uzSRbcl3N4YjO8kgemRfumHfOZz0iukwruQ",
			"p": "55-jYUpM1NLZwryeGbt0yngXwGelRtzbF4es2hyevbH257kHjy3s3Vs5VOPdlie7yB2xcl_Go5YvcYJEMQ_hn1NNac5nIpyK0adrtwzwMfKev4nDxZbxKoIoWqjqOSx1VbsO94XBjOI9qZVoWLrTCTa0csraiZSbGfnCmRux5ok",
			"q": "wYcbu7ehOUA8G0r5Q1ZFB2cu7xO-6mn1b9gO_eyIC4KlrU1IE1gVniM-81RdGxcrOrSLKakwxCbW0YsG8mJg1ZiAkqqQDx7O9Lfgd0LFCOSODNREqia7wQQXXs32Dk__yEzpj9HLwABSIk5Hp5p8L2nxRYlxIxkeykL_cXsLz7E",
			"qi": "xsMoNcfCz9er-6_GMlx4_LclYGubqDPFl9xNlpsG0on8Cop0FiXAr9LKyCgUkcsRbExdNClTDDtqULqC8sC0V9vCUQjt1RjVJPUMgtN2JMAAMsnDgc7vxi5oCDOmxkRonl2foPUYUr7MB-J_jKlLmOHjYxYAWqHy_pl7w7r0kdM"
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
