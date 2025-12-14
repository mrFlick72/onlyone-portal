package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/middleware/security"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/web/api"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(security.SetUpOAuth2())

	api.RegisterEndpoints(router, nil)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := router.Run("0.0.0.0:8088"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
