package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)


func RegisterEndpoints(r *gin.Engine) {

	// Define a simple GET endpoint
	r.GET("/api/tags", func(c *gin.Context) {
		// Return JSON response
		c.Request.Context()
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Define a simple GET endpoint
	r.PUT("/api/tags", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
