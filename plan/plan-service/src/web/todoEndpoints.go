package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/pkg/clock"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/plan"
)

var logger = logging.GetLoggerInstanceForComponentByTypeName("TodoEndpoints")

type todoRepresentation struct {
	Id       string `json:"id"`
	UserName string `json:"user_name"`
	Date     string `json:"date"`
	Content  string `json:"content"`
}

type TodoEndpoints struct {
	TodoRepository plan.TodoRepository
	ContextFactory server.ContextFactoryConverter
}

func RegisterEndpoints(r *gin.Engine, factory server.ContextFactoryConverter, repo plan.TodoRepository) {
	endpoint := &TodoEndpoints{TodoRepository: repo, ContextFactory: factory}
	contextPath := "/todo-service"
	r.GET(contextPath+"/todo", endpoint.GetTodoEndpoint)
	r.GET(contextPath+"/todo/:id", endpoint.GetOneTodoEndpoint)
	r.POST(contextPath+"/todo", endpoint.SaveTodoEndpoint)
	r.DELETE(contextPath+"/todo/:id", endpoint.DeleteTodoEndpoint)
}

func (endpoints *TodoEndpoints) GetTodoEndpoint(c *gin.Context) {
	ctx := endpoints.ContextFactory.CreateContextFromGin(c)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	allTodo, err := endpoints.TodoRepository.GetAllTodo(*user.UserName)
	if err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, fromDomainToRepresentationForAllTodoInList(allTodo))
}

func (endpoints *TodoEndpoints) GetOneTodoEndpoint(c *gin.Context) {
	id := c.Param("id")
	todo, err := endpoints.TodoRepository.GetTodo(id)
	if err != nil || todo == nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, fromDomainToRepresentation(todo))
}

func (endpoints *TodoEndpoints) SaveTodoEndpoint(c *gin.Context) {
	var rep todoRepresentation
	if err := c.ShouldBindJSON(&rep); err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusBadRequest)
		return
	}
	rep.Id = uuid.NewString()

	if err := endpoints.TodoRepository.SaveTodo(fromRepresentationToDomain(&rep)); err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusCreated)
}

func (endpoints *TodoEndpoints) DeleteTodoEndpoint(c *gin.Context) {
	id := c.Param("id")
	if err := endpoints.TodoRepository.RemoveTodo(id); err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

func fromDomainToRepresentationForAllTodoInList(allTodo []*plan.Todo) []todoRepresentation {
	result := make([]todoRepresentation, 0, len(allTodo))
	for _, todo := range allTodo {
		result = append(result, fromDomainToRepresentation(todo))
	}
	return result
}

func fromDomainToRepresentation(todo *plan.Todo) todoRepresentation {
	return todoRepresentation{
		Id:       todo.Id,
		UserName: todo.UserName,
		Date:     clock.FormatDateFor(todo.Date),
		Content:  todo.Content,
	}
}

func fromRepresentationToDomain(rep *todoRepresentation) *plan.Todo {
	return &plan.Todo{
		Id:       rep.Id,
		UserName: rep.UserName,
		Date:     clock.ParseDateFor(rep.Date),
		Content:  rep.Content,
	}
}
