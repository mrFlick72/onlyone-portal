package management

import "github.com/gin-gonic/gin"

func RegisterEndpoints(r *gin.Engine) *gin.Engine {

	// GET /api/tags — return all tags as JSON
	r.GET("/management/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	return r
}
