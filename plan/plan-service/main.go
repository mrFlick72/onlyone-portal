package main

import (
	"fmt"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/configuration"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/pkg/config"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/pkg/logging"
)

func main() {
	var configurationManager = config.GetConfigurationManagerInstance()

	logging.GetLoggerInstance()

	server := configuration.ServerConfigurer()
	fmt.Println(configurationManager.GetConfigFor("LOGGING_FILE_NAME"))
	webServerPort := configurationManager.GetConfigFor("WEB_SERVER_PORT")
	server.Start(fmt.Sprintf(":%v", webServerPort))

}
