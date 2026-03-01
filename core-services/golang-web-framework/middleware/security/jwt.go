package security

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"go.xrfang.cn/wild"
)

var configurationManager = config.GetConfigurationManagerInstance()

func SetUpOAuth2() gin.HandlerFunc {
	jwk := Jwk{
		Url: configurationManager.GetConfigFor("idp.jwks-endpoint"),
	}
	role := configurationManager.GetConfigFor("user.required-role")
	sets, _ := jwk.JwkSets()
	jwt_logger.LogInfofFor("OAuth2 middleware set up with role: %s", role)
	return NewOAuth2Middleware(sets, role, []string{"/management/*"})
}

func NewOAuth2Middleware(keySet jwk.Set, allowedAuthority string, ignored []string) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		for _, path := range ignored {
			m := wild.MustCompile(path, wild.Extended())

			if m.Match(ctx.FullPath()) {
				jwt_logger.LogInfofFor("skipping oauth2 evaluation for path: %s\n", path)
				ctx.Next()
				return
			}
		}
		if ctx.Request.Method == "OPTIONS" {
			jwt_logger.LogInfofFor("skipping oauth2 evaluation for OPTIONS request\n")
			ctx.Next()
			return
		}

		authorization := authorizationHeaderFor(ctx)
		jwt_logger.LogDebugfFor("verifying token: %s\n", authorization)
		jwt, err := jwt.Parse([]byte(authorization), jwt.WithVerify(false))
		if err != nil {
			jwt_logger.LogErrorfFor("failed to parse jwt token: %v\n", err)
			ctx.Status(401)
			ctx.Abort()
			return
		}

		exp, ok := jwt.Expiration()
		if !ok {
			jwt_logger.LogErrorfFor("failed to fetch exp from token: %v\n", err)
			ctx.Status(401)
			ctx.Abort()
			return
		}

		if time.Now().After(exp) {
			jwt_logger.LogErrorfFor("token is expired: %v\n", err)
			ctx.Status(401)
			ctx.Abort()
			return
		}

		userName := getClaimFromToken(jwt, "user_name")
		jwt_logger.LogInfofFor("authenticating user: %s\n", *userName)

		authorities := getClaimListFromToken(jwt, "authorities")
		jwt_logger.LogInfofFor("user %s has those authority %s\n", *userName, *authorities)

		if ok := contains(toStringSlice(*authorities), allowedAuthority); !ok {
			jwt_logger.LogErrorfFor("user %s does not have required authority: %s\n", *userName, allowedAuthority)
			ctx.Status(403)
			ctx.Abort()
			return
		} else {
			jwt_logger.LogInfofFor("user %s has required authority: %s\nThe user has %s", *userName, allowedAuthority, *authorities)
		}

		ctx.Set("user", User{
			UserName:    userName,
			Authorities: toStringSlice(*authorities),
			AccessToken: &authorization,
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

func authorizationHeaderFor(ctx *gin.Context) string {
	authorization := ctx.GetHeader("Authorization")
	authorization = authorization[7:] // remove "Bearer "
	return authorization
}
