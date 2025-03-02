package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konrad2002/tmate-server/dto"
	"github.com/konrad2002/tmate-server/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return UserController{
		userService: userService,
	}
}

func (uc *UserController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/user/")

	router.GET("", uc.getAllUsers)
	router.GET("id/:id", uc.getUserById)
	router.GET("username/:username", uc.getUserByUsername)

	router.POST("", uc.createUser)
	router.POST("login", uc.login)

	router.DELETE("id/:id", uc.removeUser)

	router.OPTIONS("", uc.ok)
	router.OPTIONS("login", uc.ok)

}

func (uc *UserController) ok(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (uc *UserController) getAllUsers(c *gin.Context) {
	users, err := uc.userService.GetAll()
	if err != nil {
		fmt.Printf(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func (uc *UserController) getUserById(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	user, err := uc.userService.GetUserById(id)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (uc *UserController) getUserByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no username given"})
		return
	}

	user, err := uc.userService.GetUserByUsername(username)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (uc *UserController) createUser(c *gin.Context) {
	var user dto.CreateUserDto
	if err := c.BindJSON(&user); err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := uc.userService.CreateUser(user)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, r)
}

func (uc *UserController) login(c *gin.Context) {
	var login dto.LoginDto
	if err := c.BindJSON(&login); err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := uc.userService.Login(login)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, r)
}

func (uc *UserController) removeUser(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	user, err := uc.userService.RemoveUser(id)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}
