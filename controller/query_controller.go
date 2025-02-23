package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konrad2002/tmate-server/service"
	"net/http"
)

type QueryController struct {
	queryService service.QueryService
}

func NewQueryController(queryService service.QueryService) QueryController {
	return QueryController{
		queryService: queryService,
	}
}

func (qc *QueryController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/query/")
	router.GET("", qc.getAllQueries)

	router.POST("save-example", qc.saveExample)
}

func (qc *QueryController) getAllQueries(c *gin.Context) {
	queries, err := qc.queryService.GetAll()
	if err != nil {
		fmt.Printf(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, queries)
}

func (qc *QueryController) saveExample(c *gin.Context) {
	queries, err := qc.queryService.SaveExample()
	if err != nil {
		fmt.Printf(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, queries)
}
