package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/domain"
)

func RegisterEndpoints(r *gin.Engine, repository *domain.TagRepository) *gin.Engine {

	// GET /api/tags — return all tags as JSON
	r.GET("/api/tags", func(c *gin.Context) {
		ctx := c.Request.Context()
		tags, err := (*repository).FindAllTags(&ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Return a bare JSON array (not wrapped) — empty array when there are no tags
		if tags == nil {
			c.JSON(http.StatusOK, []domain.Tag{})
			return
		}
		c.JSON(http.StatusOK, *tags)
	})

	// PUT /api/tags — create or update a tag from JSON body
	r.PUT("/api/tags", func(c *gin.Context) {
		var tag domain.Tag
		if err := c.ShouldBindJSON(&tag); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		if err := (*repository).SaveTag(&ctx, &tag); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, tag)
	})

	return r
}
