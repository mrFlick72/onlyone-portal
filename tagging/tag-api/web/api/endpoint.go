package api

import (
	"errors"
	"net/http"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/domain"
)

func RegisterEndpoints(r *gin.Engine, contextFactoryConverter server.ContextFactoryConverter, repository domain.TagRepository, findAllTagsAction *domain.FindAllTags) *gin.Engine {

	var logger = logging.GetLoggerInstanceForComponentByTypeName("api.RegisterEndpoints")

	// GET /api/tags/scope/:scope — return tags matching the given scope, including the UNKNOWN sentinel tag.
	// Scope is authoritative: there is no unscoped "return everything" route — see
	// docs/adr/0007-scope-mandatory-and-scoped-reads-are-strict.md.
	r.GET("/api/tags/scope/:scope", func(c *gin.Context) {
		tags, err := findAllTagsAction.Execute(contextFactoryConverter.CreateContextFromGin(c), c.Param("scope"))
		if err != nil {
			logger.LogErrorfFor("error occurred: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tags)
	})

	// PUT /api/tags — create a tag from JSON body. Scope is mandatory and
	// non-blank: a tag whose normalized scope is empty is rejected with 400.
	// See docs/adr/0007-scope-mandatory-and-scoped-reads-are-strict.md.
	r.PUT("/api/tags", func(c *gin.Context) {
		var tag domain.Tag
		if err := c.ShouldBindJSON(&tag); err != nil {
			logger.LogErrorfFor("error occurred: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if domain.NormalizeScope(tag.Scope) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "scope is mandatory and must not be blank"})
			return
		}

		tag.Key = uuid.New().String() // Assign a new UUID as the key

		ctx := contextFactoryConverter.CreateContextFromGin(c)
		if err := repository.SaveTag(ctx, &tag); err != nil {
			logger.LogErrorfFor("error occurred: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, tag)
	})

	// PATCH /api/tags/scope/:scope/:key — update only the Value of the tag at
	// that (key, scope) composite identity. Key and Scope are immutable; a
	// :key/:scope pair that doesn't match an existing tag is 404, and the
	// UNKNOWN sentinel key is rejected with 400 since it is never persisted.
	// See docs/adr/0008-tag-update-is-value-only.md.
	r.PATCH("/api/tags/scope/:scope/:key", func(c *gin.Context) {
		key := c.Param("key")
		scope := domain.NormalizeScope(c.Param("scope"))

		if key == domain.UnknownSentinelKey {
			c.JSON(http.StatusBadRequest, gin.H{"error": "the UNKNOWN sentinel tag cannot be updated"})
			return
		}

		var body struct {
			Value string `json:"value"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			logger.LogErrorfFor("error occurred: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := contextFactoryConverter.CreateContextFromGin(c)
		err := repository.UpdateTagValue(ctx, &domain.Tag{Key: key, Scope: scope, Value: body.Value})
		if errors.Is(err, domain.ErrTagNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err != nil {
			logger.LogErrorfFor("error occurred: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	})

	// DELETE /api/tags/scope/:scope/:key — delete the tag at that (key, scope)
	// composite identity. Same 404/400 rules as PATCH above.
	r.DELETE("/api/tags/scope/:scope/:key", func(c *gin.Context) {
		key := c.Param("key")
		scope := domain.NormalizeScope(c.Param("scope"))

		if key == domain.UnknownSentinelKey {
			c.JSON(http.StatusBadRequest, gin.H{"error": "the UNKNOWN sentinel tag cannot be deleted"})
			return
		}

		ctx := contextFactoryConverter.CreateContextFromGin(c)
		err := repository.DeleteTag(ctx, key, scope)
		if errors.Is(err, domain.ErrTagNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err != nil {
			logger.LogErrorfFor("error occurred: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	})

	return r
}
