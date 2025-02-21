package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konrad2002/tmate-server/model"
	"github.com/konrad2002/tmate-server/service"
	"net/http"
)

type FieldController struct {
	fieldService service.FieldService
}

func NewFieldController(fieldService service.FieldService) FieldController {
	return FieldController{
		fieldService: fieldService,
	}
}

func (mc *FieldController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/field")
	router.GET("", mc.getAllFields)
	router.GET("/types", mc.getFieldTypes)
}

func (mc *FieldController) getAllFields(c *gin.Context) {
	fields, err := mc.fieldService.GetAll()
	if err != nil {
		fmt.Printf(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, fields)
}

func (mc *FieldController) getFieldTypes(c *gin.Context) {
	fieldTypes := model.GetAllFieldType()

	c.IndentedJSON(http.StatusOK, fieldTypes)
}
