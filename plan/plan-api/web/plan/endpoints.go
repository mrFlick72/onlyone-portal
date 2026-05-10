package plan

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/domain/plan"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/pkg/clock"
)

var logger = logging.GetLoggerInstanceForComponentByTypeName("PlanEndpoints")

type PlanEndpoints struct {
	repository plan.PlanRepository
	factory    server.ContextFactoryConverter
}

func RegisterEndpoints(r *gin.Engine, factory server.ContextFactoryConverter, repo plan.PlanRepository) {
	e := &PlanEndpoints{repository: repo, factory: factory}
	g := r.Group("/api")
	g.GET("/plan", e.getAll)
	g.GET("/plan/:id", e.getOne)
	g.POST("/plan", e.save)
	g.PUT("/plan/:id/todo", e.addTodo)
	g.PUT("/plan/:id/todo/:todoId", e.updateTodo)
	g.PUT("/plan/:id/todo/:todoId/status", e.changeTodoStatus)
	g.DELETE("/plan/:id/todo/:todoId", e.removeTodo)
	g.DELETE("/plan/:id", e.delete)
}

func (e *PlanEndpoints) getAll(c *gin.Context) {
	ctx := e.factory.CreateContextFromGin(c)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	plans, err := e.repository.GetAllPlanBy(*user.UserName)
	if err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, toPlanRepresentationList(plans))
}

func (e *PlanEndpoints) getOne(c *gin.Context) {
	ctx := e.factory.CreateContextFromGin(c)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	p, err := e.repository.GetPlan(c.Param("id"), *user.UserName)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, toPlanRepresentation(p))
}

func (e *PlanEndpoints) save(c *gin.Context) {
	ctx := e.factory.CreateContextFromGin(c)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	var rep planRepresentation
	if err := c.ShouldBindJSON(&rep); err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusBadRequest)
		return
	}
	planId, err := e.repository.CreateNewPlan(toPlanDomain(&rep, *user.UserName))
	if err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": planId})
}

func (e *PlanEndpoints) addTodo(c *gin.Context) {
	ctx := e.factory.CreateContextFromGin(c)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	var rep todoRepresentation
	if err := c.ShouldBindJSON(&rep); err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusBadRequest)
		return
	}
	t := plan.Todo{
		Id:       uuid.NewString(),
		UserName: *user.UserName,
		Date:     clock.ParseDateFor(rep.Date),
		Content:  rep.Content,
		Status:   plan.StatusTodo,
	}
	if err := e.repository.AddTodo(c.Param("id"), t); err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusCreated)
}

func (e *PlanEndpoints) removeTodo(c *gin.Context) {
	if err := e.repository.RemoveTodo(c.Param("id"), c.Param("todoId")); err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

func (e *PlanEndpoints) updateTodo(c *gin.Context) {
	ctx := e.factory.CreateContextFromGin(c)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	var rep todoRepresentation
	if err := c.ShouldBindJSON(&rep); err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusBadRequest)
		return
	}
	t := plan.Todo{
		Id:       c.Param("todoId"),
		UserName: *user.UserName,
		Date:     clock.ParseDateFor(rep.Date),
		Content:  rep.Content,
	}
	if err := e.repository.UpdateTodo(c.Param("id"), t); err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

func (e *PlanEndpoints) changeTodoStatus(c *gin.Context) {
	ctx := e.factory.CreateContextFromGin(c)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	var rep todoStatusChangeRequest
	if err := c.ShouldBindJSON(&rep); err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusBadRequest)
		return
	}
	target := plan.TodoStatus(rep.Status)
	if !target.IsValid() {
		c.Status(http.StatusBadRequest)
		return
	}
	planId := c.Param("id")
	todoId := c.Param("todoId")
	p, err := e.repository.GetPlan(planId, *user.UserName)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	var current *plan.Todo
	for _, t := range p.Todos {
		if t.Id == todoId {
			current = t
			break
		}
	}
	if current == nil {
		c.Status(http.StatusNotFound)
		return
	}
	if !current.Status.CanTransitionTo(target) {
		c.Status(http.StatusConflict)
		return
	}
	if err := e.repository.UpdateTodoStatus(planId, todoId, target); err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

func (e *PlanEndpoints) delete(c *gin.Context) {
	ctx := e.factory.CreateContextFromGin(c)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	if err := e.repository.DeletePlan(c.Param("id"), *user.UserName); err != nil {
		logger.LogErrorFor(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

func toPlanRepresentationList(plans []*plan.Plan) []planRepresentation {
	result := make([]planRepresentation, 0, len(plans))
	for _, p := range plans {
		result = append(result, toPlanRepresentation(p))
	}
	return result
}

func toPlanRepresentation(p *plan.Plan) planRepresentation {
	todos := make([]todoRepresentation, 0, len(p.Todos))
	for _, t := range p.Todos {
		todos = append(todos, todoRepresentation{
			Id:       t.Id,
			UserName: t.UserName,
			Date:     clock.FormatDateFor(t.Date),
			Content:  t.Content,
			Status:   string(t.Status),
		})
	}
	return planRepresentation{
		Id:        p.Id,
		UserName:  p.UserName,
		Title:     p.Title,
		Date:      clock.FormatDateFor(p.Date),
		Todos:     todos,
		TodoCount: len(todos),
	}
}

func toPlanDomain(rep *planRepresentation, userName string) plan.Plan {
	return plan.Plan{
		UserName: userName,
		Title:    rep.Title,
		Date:     clock.ParseDateFor(rep.Date),
	}
}
