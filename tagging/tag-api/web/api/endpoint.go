package api

import (
	"log"
	"net/http"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/domain"
)

func RegisterEndpoints(r *gin.Engine, contextFactoryConverter server.ContextFactoryConverter, repository domain.TagRepository) *gin.Engine {

	// GET /api/tags — return all tags as JSON
	r.GET("/api/tags", func(c *gin.Context) {
		tags, err := repository.FindAllTags(contextFactoryConverter.CreateContextFromGin(c))
		if err != nil {
			log.Println("error occurred:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Return a bare JSON array (not wrapped) — empty array when there are no tags
		if tags == nil {
			c.JSON(http.StatusOK, []domain.Tag{})
			return
		}
		c.JSON(http.StatusOK, tags)
	})

	// PUT /api/tags — create or update a tag from JSON body
	r.PUT("/api/tags", func(c *gin.Context) {
		var tag domain.Tag
		if err := c.ShouldBindJSON(&tag); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tag.Key = uuid.New().String() // Assign a new UUID as the key
		log.Println("Saving tag:", tag)

		ctx := contextFactoryConverter.CreateContextFromGin(c)
		if err := repository.SaveTag(ctx, &tag); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, tag)
	})

	return r
}
