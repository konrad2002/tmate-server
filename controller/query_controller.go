package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konrad2002/tmate-server/model"
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
	router.POST("", qc.addQuery)

	router.PUT("", qc.updateQuery)

	router.OPTIONS("", qc.ok)

}

func (qc *QueryController) ok(c *gin.Context) {
	c.Status(http.StatusOK)
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

func (qc *QueryController) addQuery(c *gin.Context) {
	var query model.Query
	if err := c.BindJSON(&query); err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := qc.queryService.AddQuery(query)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func (qc *QueryController) updateQuery(c *gin.Context) {
	var query model.Query
	if err := c.BindJSON(&query); err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := qc.queryService.UpdateQuery(query)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}
