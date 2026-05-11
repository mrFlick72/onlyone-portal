package i18n-api

import (
"github.com/lestrrat-go/jwx/v3/jwt"
)
type LanguageResolver struct {
	ClaimName       string
	DefaultLanguage Locale
}

func NewLanguageResolver(claimName string, defaultLanguage Locale) *LanguageResolver {
	return &LanguageResolver{ClaimName: claimName, DefaultLanguage: defaultLanguage.Normalized()}
}

// ResolveFromAccessToken parses the bearer token (already verified by the
// OAuth2 middleware) and returns the locale stored in the configured claim.
// Falls back to the default language when the token cannot be parsed or the
// claim is missing/empty.
func (r *LanguageResolver) ResolveFromAccessToken(accessToken string) Locale {
	if accessToken == "" {
		return r.DefaultLanguage
	}
	token, err := jwt.ParseInsecure([]byte(accessToken))
	if err != nil {
		return r.DefaultLanguage
	}
	var claim string
	if err := token.Get(r.ClaimName, &claim); err != nil || claim == "" {
		return r.DefaultLanguage
	}
	return Locale(claim).Normalized()
}
