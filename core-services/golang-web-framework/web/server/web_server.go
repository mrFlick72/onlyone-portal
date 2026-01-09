package server

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/magangement"

	"github.com/gin-gonic/gin"
)

var configurationManager = config.GetConfigurationManagerInstance()
 
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
		AllowOrigins:     strings.Split(configurationManager.GetConfigFor("CORS_ALLOWED_ORIGINS"), ","),
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin","Authorization", "Content-Type", "Accept"},
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
	port := os.Getenv("WEBSERVER_PORT")
	serverBinder := fmt.Sprintf("0.0.0.0:%s", port)
	if err := wsp.router.Run(serverBinder); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
	return nil
}
