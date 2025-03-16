package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konrad2002/tmate-server/auth"
	"github.com/konrad2002/tmate-server/dto"
	"github.com/konrad2002/tmate-server/model"
	"github.com/konrad2002/tmate-server/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type QueryController struct {
	queryService service.QueryService
	userService  service.UserService
}

func NewQueryController(queryService service.QueryService, userService service.UserService) QueryController {
	return QueryController{
		queryService: queryService,
		userService:  userService,
	}
}

func (qc *QueryController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/query/")

	router.Use(auth.HandlerFunc(&qc.userService))

	router.GET("", qc.getAllQueries)
	router.GET("me", qc.getAllQueriesForMe)
	router.GET("id/:id", qc.getQueryById)

	router.POST("save-example", qc.saveExample)
	router.POST("", qc.addQuery)

	router.PUT("", qc.updateQuery)

	router.DELETE(":id", qc.removeQuery)

	router.OPTIONS("", qc.ok)
	router.OPTIONS(":id", qc.ok)

}

func (qc *QueryController) ok(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (qc *QueryController) getAllQueries(c *gin.Context) {
	queries, err := qc.queryService.GetAll()
	if err != nil {
		fmt.Print(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, queries)
}

func (qc *QueryController) getAllQueriesForMe(c *gin.Context) {
	u, _ := c.Get("currentUser")
	user := u.(dto.UserInfoDto)
	queries, err := qc.queryService.GetAllForUser(user.Identifier)
	if err != nil {
		fmt.Print(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, queries)
}

func (qc *QueryController) getQueryById(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	query, err := qc.queryService.GetQueryById(id)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, query)
}

func (qc *QueryController) saveExample(c *gin.Context) {
	queries, err := qc.queryService.SaveExample()
	if err != nil {
		fmt.Print(err.Error())
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

	// set owner id to submitting person if non given or user has no permission to submit own
	u, _ := c.Get("currentUser")
	user := u.(dto.UserInfoDto)

	if !user.Permissions.QueryManagement {
		query.OwnerUserId = user.Identifier
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

	// set owner id to submitting person if non given or user has no permission to submit own
	u, _ := c.Get("currentUser")
	user := u.(dto.UserInfoDto)

	if !user.Permissions.QueryManagement {
		query.OwnerUserId = user.Identifier
	}

	r, err := qc.queryService.UpdateQuery(query)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func (qc *QueryController) removeQuery(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	query, err := qc.queryService.RemoveQuery(id)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, query)
}
