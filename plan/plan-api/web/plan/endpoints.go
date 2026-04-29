package plan

import (
	"github.com/gin-gonic/gin"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
)

var logger = logging.GetLoggerInstanceForComponentByTypeName("TodoEndpoints")

type todoRepresentation struct {
	Id       string `json:"id"`
	UserName string `json:"user_name"`
	Date     string `json:"date"`
	Content  string `json:"content"`
}

type TodoEndpoints struct {
	repository plan.PlanRepository
	factory    server.ContextFactoryConverter
}

func RegisterEndpoints(r *gin.Engine, factory server.ContextFactoryConverter, repo plan.PlanRepository) {
	e := &TodoEndpoints{repository: repo, factory: factory}
	g := r.Group("/api")
	g.GET("/plan", e.getAll)
	g.GET("/plan/:id", e.getOne)
	g.POST("/plan", e.save)
	g.PUT("/plan/:id/todo", e.addTodo)
	g.PUT("/plan/:id/todo/:todoId", e.updateTodo)
	g.DELETE("/plan/:id/todo/:todoId", e.removeTodo)
	g.DELETE("/plan/:id", e.delete)
}

func (e *TodoEndpoints) getAll(c *gin.Context) {
	//ctx := e.factory.CreateContextFromGin(c)
	//user, err := security.GetCurrentUser(ctx)
	//if err != nil {
	//	c.Status(http.StatusUnauthorized)
	//	return
	//}
	//
	//todos, err := e.repository.GetAllTodo(*user.UserName)
	//if err != nil {
	//	logger.LogErrorFor(err)
	//	c.Status(http.StatusInternalServerError)
	//	return
	//}
	//c.JSON(http.StatusOK, toRepresentationList(todos))
}

func (e *TodoEndpoints) getOne(c *gin.Context) {
	//t, err := e.repository.GetTodo(c.Param("id"))
	//if err != nil || t == nil {
	//	c.Status(http.StatusNotFound)
	//	return
	//}
	//c.JSON(http.StatusOK, toRepresentation(t))
}

func (e *TodoEndpoints) save(c *gin.Context) {
	//var rep todoRepresentation
	//if err := c.ShouldBindJSON(&rep); err != nil {
	//	logger.LogErrorFor(err)
	//	c.Status(http.StatusBadRequest)
	//	return
	//}
	//rep.Id = uuid.NewString()
	//
	//if err := e.repository.SaveTodo(toDomain(&rep)); err != nil {
	//	logger.LogErrorFor(err)
	//	c.Status(http.StatusInternalServerError)
	//	return
	//}
	//c.Status(http.StatusCreated)
}

func (e *TodoEndpoints) delete(c *gin.Context) {
	//if err := e.repository.RemoveTodo(c.Param("id")); err != nil {
	//	logger.LogErrorFor(err)
	//	c.Status(http.StatusInternalServerError)
	//	return
	//}
	//c.Status(http.StatusNoContent)
}

func (e *TodoEndpoints) addTodo(context *gin.Context) {

}

func (e *TodoEndpoints) updateTodo(context *gin.Context) {

}

func (e *TodoEndpoints) removeTodo(context *gin.Context) {

}

//func toRepresentationList(todos []*todo.Todo) []todoRepresentation {
//	result := make([]todoRepresentation, 0, len(todos))
//	for _, t := range todos {
//		result = append(result, toRepresentation(t))
//	}
//	return result
//}

//func toRepresentation(t *todo.Todo) todoRepresentation {
//	return todoRepresentation{
//		Id:       t.Id,
//		UserName: t.UserName,
//		Date:     clock.FormatDateFor(t.Date),
//		Content:  t.Content,
//	}
//}

//func toDomain(rep *todoRepresentation) *todo.Todo {
//	return &todo.Todo{
//		Id:       rep.Id,
//		UserName: rep.UserName,
//		Date:     clock.ParseDateFor(rep.Date),
//		Content:  rep.Content,
//	}
//}
