package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrflick72/budget/budget-api/domain/budget/revenue"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

func RegisterRevenueEndpoints(
	r *gin.Engine,
	ContextFactoryConverter server.ContextFactoryConverter,
	facade revenue.RevenueActions,
) *gin.Engine {

	var logger = logging.GetLoggerInstanceForComponentByTypeName("web.RegisterRevenueEndpoints")

	r.GET("/api/budget/revenue", func(c *gin.Context) {
		year, err := parseYearQueryParam(c.Query("q"))
		if err != nil {
			logger.LogErrorfFor("Error parsing revenue query param: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := ContextFactoryConverter.CreateContextFromGin(c)
		revenues, err := facade.FindRevenue(ctx, year)
		if err != nil {
			logger.LogErrorfFor("Error finding revenues: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, RevenueListDomainToRepresentationModel(revenues))
	})

	r.POST("/api/budget/revenue", func(c *gin.Context) {
		var representation RevenueRepresentation

		ctx := ContextFactoryConverter.CreateContextFromGin(c)
		if err := c.ShouldBindJSON(&representation); err != nil {
			logger.LogErrorfFor("Error binding JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		domainModel := RevenueRepresentationToDomainModel(representation)
		facade.CreateRevenue(ctx, domainModel)
		c.Status(http.StatusCreated)
	})

	r.PUT("/api/budget/revenue/:id", func(c *gin.Context) {
		var representation RevenueRepresentation

		ctx := ContextFactoryConverter.CreateContextFromGin(c)
		if err := c.ShouldBindJSON(&representation); err != nil {
			logger.LogErrorfFor("Error binding JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		domainModel := RevenueRepresentationToDomainModel(representation)
		domainModel.Id = c.Param("id")
		facade.UpdateRevenue(ctx, domainModel)
		c.Status(http.StatusNoContent)
	})

	r.DELETE("/api/budget/revenue/:id", func(c *gin.Context) {
		ctx := ContextFactoryConverter.CreateContextFromGin(c)
		facade.DeleteRevenue(ctx, c.Param("id"))
		c.Status(http.StatusNoContent)
	})

	return r
}
