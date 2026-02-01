package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konrad2002/tmate-server/auth"
	"github.com/konrad2002/tmate-server/dto"
	"github.com/konrad2002/tmate-server/model"
	"github.com/konrad2002/tmate-server/service"
	"net/http"
)

type FieldController struct {
	fieldService service.FieldService
	userService  service.UserService
}

func NewFieldController(fieldService service.FieldService, userService service.UserService) FieldController {
	return FieldController{
		fieldService: fieldService,
		userService:  userService,
	}
}

func (fc *FieldController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/field/")

	router.Use(auth.HandlerFunc(&fc.userService))

	router.GET("", fc.getAllFields)
	router.GET("types", fc.getFieldTypes)

	router.POST("", fc.addField)
}

func (fc *FieldController) getAllFields(c *gin.Context) {
	fields, err := fc.fieldService.GetAll()
	if err != nil {
		fmt.Print(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, fields)
}

func (fc *FieldController) getFieldTypes(c *gin.Context) {
	fieldTypes := model.GetAllFieldType()

	c.IndentedJSON(http.StatusOK, fieldTypes)
}

func (fc *FieldController) addField(c *gin.Context) {

	u, _ := c.Get("currentUser")
	user := u.(dto.UserInfoDto)
	if !user.Permissions.TableStructureManagement {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "no table structure management permissions"})
		return
	}

	var field model.Field
	if err := c.BindJSON(&field); err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := fc.fieldService.AddField(field)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}
