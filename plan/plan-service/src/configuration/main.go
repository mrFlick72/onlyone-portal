package configuration

import (
	"fmt"

	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/pkg/config"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/plan"
	api "github.com/mrflick72/onlyone-portal/plan/plan-service/src/web"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var configurationManager = config.GetConfigurationManagerInstance()

func ServerConfigurer() *echo.Echo {
	repository := plan.MySqlTodoRepository{ConnectionString: databaseConnectionStringFromEnv()}
	server := echo.New()
	contextPath := "/todo-service"
	endpoint := api.TodoEndpoints{TodoRepository: &repository}

	server.GET(contextPath+"/todo", endpoint.GetTodoEndpoint)
	server.GET(contextPath+"/todo/:id", endpoint.GetOneTodoEndpoint)
	server.POST(contextPath+"/todo", endpoint.SaveTodoEndpoint)
	server.DELETE(contextPath+"/todo/:id", endpoint.DeleteTodoEndpoint)

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "DELETE"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	server.Use(middleware.Logger())

	return server
}

func databaseConnectionStringFromEnv() string {
	databaseUrl := configurationManager.GetConfigFor("DATABASE_URL")
	databaseUserName := configurationManager.GetConfigFor("DATABASE_USER")
	databasePassword := configurationManager.GetConfigFor("DATABASE_PASSWORD")
	databaseConnectionString := fmt.Sprintf("%v:%v@tcp(%v)/todo?parseTime=true", databaseUserName, databasePassword, databaseUrl)
	return databaseConnectionString
}
