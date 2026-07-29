package security

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jws"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"go.xrfang.cn/wild"
)

var configurationManager = config.GetConfigurationManagerInstance()

// SetUpOAuth2 fetches the IDP JWKS (with periodic background refresh bound to
// ctx) and returns a Gin middleware that verifies bearer tokens against it.
func SetUpOAuth2(ctx context.Context) (gin.HandlerFunc, error) {
	jwksURL := configurationManager.GetConfigFor("idp.jwks-endpoint")
	role := configurationManager.GetConfigFor("user.required-role")

	keySet, err := NewCachedJwkSet(ctx, jwksURL)
	if err != nil {
		jwt_logger.LogErrorfFor("oauth2 setup: %v", err)
		return nil, fmt.Errorf("oauth2 setup: %w", err)
	}

	jwt_logger.LogInfofFor("OAuth2 middleware set up with role: %s", role)
	return NewOAuth2Middleware(keySet, role, []string{"/management/*"}), nil
}

func NewOAuth2Middleware(keySet jwk.Set, allowedAuthority string, ignored []string) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		for _, path := range ignored {
			m := wild.MustCompile(path, wild.Extended())

			if m.Match(ctx.FullPath()) {
				jwt_logger.LogDebugfFor("skipping oauth2 evaluation for path: %s\n", path)
				ctx.Next()
				return
			}
		}
		if ctx.Request.Method == "OPTIONS" {
			jwt_logger.LogDebugfFor("skipping oauth2 evaluation for OPTIONS request\n")
			ctx.Next()
			return
		}

		bearer, ok := bearerTokenFor(ctx)
		if !ok {
			jwt_logger.LogErrorFor("missing or malformed Authorization header")
			ctx.Status(401)
			ctx.Abort()
			return
		}

		token, err := jwt.Parse([]byte(bearer),
			jwt.WithKeySet(keySet, jws.WithInferAlgorithmFromKey(true)),
		)
		if err != nil {
			jwt_logger.LogErrorfFor("failed to verify jwt token: %v\n", err)
			ctx.Status(401)
			ctx.Abort()
			return
		}

		userName := getClaimFromToken(token, "user_name")
		if userName == nil {
			jwt_logger.LogErrorFor("token missing required claim: user_name")
			ctx.Status(401)
			ctx.Abort()
			return
		}
		jwt_logger.LogInfofFor("authenticating user: %s\n", *userName)

		authorities := getClaimListFromToken(token, "authorities")
		if authorities == nil {
			jwt_logger.LogErrorfFor("token missing required claim: authorities (user %s)", *userName)
			ctx.Status(403)
			ctx.Abort()
			return
		}
		jwt_logger.LogInfofFor("user %s has those authority %s\n", *userName, *authorities)

		if ok := contains(toStringSlice(*authorities), allowedAuthority); !ok {
			jwt_logger.LogErrorfFor("user %s does not have required authority: %s\n", *userName, allowedAuthority)
			ctx.Status(403)
			ctx.Abort()
			return
		}
		jwt_logger.LogInfofFor("user %s has required authority: %s", *userName, allowedAuthority)

		ctx.Set("user", User{
			UserName:    userName,
			Authorities: toStringSlice(*authorities),
			AccessToken: &bearer,
		})

		jwt_logger.LogInfofFor("authenticated user: %s\n", *userName)
		ctx.Next()
	}

}

func getClaimFromToken(token jwt.Token, claimName string) *string {
	var claim string
	if err := token.Get(claimName, &claim); err != nil {
		jwt_logger.LogErrorfFor("failed to fetch private claim: %v\n", err)
		return nil
	}
	return &claim
}

func getClaimListFromToken(token jwt.Token, claimName string) *[]string {
	var raw []interface{}
	if err := token.Get(claimName, &raw); err != nil {
		jwt_logger.LogErrorfFor("failed to fetch claim authorities: %v\n", err)
		return nil
	}

	roles := make([]string, 0, len(raw))
	for _, v := range raw {
		s, ok := v.(string)
		if !ok {
			jwt_logger.LogErrorfFor("element is not a string")
			return nil
		}
		roles = append(roles, s)
	}
	return &roles
}

func contains(slice *[]string, item string) bool {
	for _, s := range *slice {
		if s == item {
			return true
		}
	}
	return false
}
func toStringSlice(slice []string) *[]string {
	result := make([]string, 0)

	for _, item := range slice {
		result = append(result, item)
	}

	return &result
}

// bearerTokenFor extracts the token from an "Authorization: Bearer <token>"
// header. Returns (token, true) when the prefix is present and the token is
// non-empty, otherwise ("", false). Never panics on malformed input.
func bearerTokenFor(ctx *gin.Context) (string, bool) {
	header := ctx.GetHeader("Authorization")
	token, ok := strings.CutPrefix(header, "Bearer ")
	if !ok || token == "" {
		return "", false
	}
	return token, true
}
