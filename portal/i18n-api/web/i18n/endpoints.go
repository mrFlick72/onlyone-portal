package i18n-api

import (
"errors"
"net/http"

"github.com/gin-gonic/gin"
"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
"github.com/mrflick72/onlyone-portal/portal/i18n-api/domain/i18n-api"
)

var logger = logging.GetLoggerInstanceForComponentByTypeName("i18n-apiEndpoints")

type Endpoints struct {
	repository i18n-api.BundleRepository
	resolver   *i18n-api.LanguageResolver
	factory    server.ContextFactoryConverter
}

func RegisterEndpoints(
	r *gin.Engine,
	factory server.ContextFactoryConverter,
	repository i18n-api.BundleRepository,
	resolver *i18n-api.LanguageResolver,
) {
	e := &Endpoints{repository: repository, resolver: resolver, factory: factory}
	g := r.Group("/api")
	g.GET("/i18n-api/:application/:page", e.getBundle)
}

func (e *Endpoints) getBundle(c *gin.Context) {
	ctx := e.factory.CreateContextFromGin(c)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	language := i18n - api.Locale(c.Query("lang")).Normalized()
	if language.IsEmpty() && user.AccessToken != nil {
		language = e.resolver.ResolveFromAccessToken(*user.AccessToken)
	}
	if language.IsEmpty() {
		language = e.resolver.DefaultLanguage
	}

	key := i18n - api.BundleKey{
		Application: c.Param("application"),
		Page:        c.Param("page"),
		Language:    language,
	}

	bundle, err := e.repository.GetBundle(ctx, key)
	if errors.Is(err, i18n-api.ErrBundleNotFound) && language != e.resolver.DefaultLanguage {
		logger.LogInfofFor("bundle %s/%s not found for language %s, falling back to %s",
			key.Application, key.Page, language, e.resolver.DefaultLanguage)
		key.Language = e.resolver.DefaultLanguage
		bundle, err = e.repository.GetBundle(ctx, key)
	}
	if errors.Is(err, i18n-api.ErrBundleNotFound) {
		c.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"application": bundle.Application,
		"page":        bundle.Page,
		"language":    string(bundle.Language),
		"messages":    bundle.Messages,
	})
}
