package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/domain"
)

func RegisterEndpoints(r *gin.Engine, repository domain.TagRepository) *gin.Engine {

	// GET /api/tags — return all tags as JSON
	r.GET("/api/tags", func(c *gin.Context) {
		tags, err := repository.FindAllTags(CopyGinKeysToRequestContext(c))
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
		c.JSON(http.StatusOK, *tags)
	})

	// PUT /api/tags — create or update a tag from JSON body
	r.PUT("/api/tags", func(c *gin.Context) {
		var tag domain.Tag
		if err := c.ShouldBindJSON(&tag); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := CopyGinKeysToRequestContext(c)
		if err := repository.SaveTag(ctx, &tag); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, tag)
	})

	return r
}


func CopyGinKeysToRequestContext(c *gin.Context) *context.Context {
    newCtx := c.Request.Context()
    for k, v := range c.Keys {
        newCtx = context.WithValue(newCtx,k, v)
    }
   return &newCtx
}