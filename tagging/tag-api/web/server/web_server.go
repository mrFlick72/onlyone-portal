package server

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/middleware/security"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/web/magangement"
)

type WebServerProvisioner struct {
	router *gin.Engine
}

func (wsp *WebServerProvisioner) ConfigureEngine() *gin.Engine {
	if wsp.router != nil {
		return wsp.router
	}

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		AllowCredentials: true,
		MaxAge:           60 * time.Minute,
	}))

	log.Println("Setting up OAuth2 middleware")
	router.Use(security.SetUpOAuth2())

	magangement.RegisterEndpoints(router)
	
	wsp.router = router

	return router
}

func (wsp *WebServerProvisioner) StartEngine() error {
	if err := wsp.router.Run("0.0.0.0:8000"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
	return nil
}
