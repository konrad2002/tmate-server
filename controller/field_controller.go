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

func (fc *FieldController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/field/")
	router.GET("", fc.getAllFields)
	router.GET("types", fc.getFieldTypes)
}

func (fc *FieldController) getAllFields(c *gin.Context) {
	fields, err := fc.fieldService.GetAll()
	if err != nil {
		fmt.Printf(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, fields)
}

func (fc *FieldController) getFieldTypes(c *gin.Context) {
	fieldTypes := model.GetAllFieldType()

	c.IndentedJSON(http.StatusOK, fieldTypes)
}
