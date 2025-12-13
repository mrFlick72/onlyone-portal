package security

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func SetUpOAuth2(app *gin.Engine, jwk Jwk, role string) {
	sets, _ := jwk.JwkSets()
	var middleware = NewOAuth2Middleware(sets, role)
	app.Use(middleware)
}

func NewOAuth2Middleware(keySet jwk.Set, allowedAuthority string) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		authorization := authorizationHeaderFor(ctx)
		jwt, _ := jwt.ParseString(authorization)
		iat, _ := jwt.IssuedAt()
		if time.Now().After(iat) {
			ctx.Status(401)
			return
		}

		userName := getClaimFromToken(jwt, "user_name").(string)
		authorities, _ := getClaimFromToken(jwt, "authorities").([]interface{})

		if ok := contains(*toStringSlice(authorities), allowedAuthority); !ok {
			ctx.Status(403)
			return
		}

		ctx.Set("user", OAuth2User{
			UserName: userName,
		})

		ctx.Next()
	}
}

func getClaimFromToken(token jwt.Token, claimName string) interface{} {
	var claim string
	if err := token.Get(claimName, &claim); err != nil {
		fmt.Printf("failed to fetch private claim\n")
		return nil
	}
	return claim
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
func toStringSlice(slice []interface{}) *[]string {
	result := make([]string, 0)

	for _, item := range slice {
		result = append(result, item.(string))
	}

	return &result
}

func authorizationHeaderFor(ctx *gin.Context) string {
	authorization := ctx.GetHeader("Authorization")
	authorization = authorization[7:len(authorization)]
	return authorization
}

type OAuth2User struct {
	UserName    string
	Authorities []string
}
