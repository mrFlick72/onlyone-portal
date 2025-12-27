package security

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/domain"
	"go.xrfang.cn/wild"
)

func SetUpOAuth2() gin.HandlerFunc {
	jwk := Jwk{
		Url: os.Getenv("JWKS_ENDPOINT"),
	}
	role := os.Getenv("OAUTH2_ROLE")
	sets, _ := jwk.JwkSets()
	log.Println("OAuth2 middleware set up with role:", role)
	return NewOAuth2Middleware(sets, role, []string{"/management/*"})
}

func NewOAuth2Middleware(keySet jwk.Set, allowedAuthority string, ignored []string) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		for _, path := range ignored {
			m := wild.MustCompile(path, wild.Extended())

			if m.Match(ctx.FullPath()) {
				log.Printf("skipping oauth2 evaluation for path: %s\n", path)
				ctx.Next()
				return
			}
		}
		if ctx.Request.Method == "OPTIONS" {
			log.Printf("skipping oauth2 evaluation for OPTIONS request\n")
			ctx.Next()
			return
		}	
		
		authorization := authorizationHeaderFor(ctx)
		log.Printf("verifying token: %s\n", authorization)
		jwt, err := jwt.Parse([]byte(authorization), jwt.WithVerify(false))
		if err != nil {
			log.Printf("failed to parse jwt token: %v\n", err)
			ctx.Status(401)
			ctx.Abort()
			return
		}

		exp, ok := jwt.Expiration()
		if !ok {
			log.Printf("failed to fetch exp from token: %v\n", err)
			ctx.Status(401)
			ctx.Abort()
			return
		}

		if time.Now().After(exp) {
			log.Printf("token is expired: %v\n", err)
			ctx.Status(401)
			ctx.Abort()
			return
		}

		userName := getClaimFromToken(jwt, "user_name")
		log.Printf("authenticating user: %s\n", *userName)

		authorities := getClaimListFromToken(jwt, "authorities")
		log.Printf("user %s has those authority %s\n", *userName, *authorities)

		if ok := contains(toStringSlice(*authorities), allowedAuthority); !ok {
			log.Printf("user %s does not have required authority: %s\n", *userName, allowedAuthority)
			ctx.Status(403)
			ctx.Abort()
			return
		} else {
			log.Printf("user %s has required authority: %s\nThe user has %s", *userName, allowedAuthority, *authorities)
		}

		ctx.Set("user", domain.User{
			UserName:    userName,
			Authorities: toStringSlice(*authorities),
		})

		log.Printf("authenticated user: %s\n", *userName)
		ctx.Next()
	}

}

func getClaimFromToken(token jwt.Token, claimName string) *string {
	var claim string
	if err := token.Get(claimName, &claim); err != nil {
		log.Printf("failed to fetch private claim\n")
		return nil
	}
	return &claim
}

func getClaimListFromToken(token jwt.Token, claimName string) *[]string {
	var raw []interface{}
	if err := token.Get(claimName, &raw); err != nil {
		log.Println("failed to fetch claim authorities:", err)
		return nil
	}

	roles := make([]string, 0, len(raw))
	for _, v := range raw {
		s, ok := v.(string)
		if !ok {
			log.Println("element is not a string")
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
