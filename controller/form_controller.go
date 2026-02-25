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

type FormController struct {
	formService service.FormService
	userService service.UserService
}

func NewFormController(formService service.FormService, userService service.UserService) FormController {
	return FormController{
		formService: formService,
		userService: userService,
	}
}

func (fc *FormController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/form/")

	router.Use(auth.HandlerFunc(&fc.userService))

	router.GET("", fc.getAllForms)
	router.GET("id/:id", fc.getFormById)

	router.POST("", fc.addForm)

	router.PUT("", fc.updateForm)

	router.DELETE(":id", fc.removeForm)

	router.OPTIONS("", fc.ok)
	router.OPTIONS(":id", fc.ok)

}

func (fc *FormController) ok(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (fc *FormController) getAllForms(c *gin.Context) {
	forms, err := fc.formService.GetAll()
	if err != nil {
		fmt.Print(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, forms)
}

func (fc *FormController) getFormById(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	form, err := fc.formService.GetFormById(id)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, form)
}

func (fc *FormController) addForm(c *gin.Context) {
	u, _ := c.Get("currentUser")
	currentUser := u.(dto.UserInfoDto)
	if !currentUser.Permissions.FormManagement {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "no form management permissions"})
		return
	}

	var form model.Form
	if err := c.BindJSON(&form); err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	r, err := fc.formService.AddForm(form)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func (fc *FormController) updateForm(c *gin.Context) {
	u, _ := c.Get("currentUser")
	currentUser := u.(dto.UserInfoDto)
	if !currentUser.Permissions.FormManagement {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "no form management permissions"})
		return
	}

	var form model.Form
	if err := c.BindJSON(&form); err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	r, err := fc.formService.UpdateForm(form)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func (fc *FormController) removeForm(c *gin.Context) {
	u, _ := c.Get("currentUser")
	currentUser := u.(dto.UserInfoDto)
	if !currentUser.Permissions.FormManagement {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "no form management permissions"})
		return
	}

	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	form, err := fc.formService.RemoveForm(id)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, form)
}
