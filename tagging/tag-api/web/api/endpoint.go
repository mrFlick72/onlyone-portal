package api

import (
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

	return r
}
