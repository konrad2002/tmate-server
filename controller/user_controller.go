package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konrad2002/tmate-server/auth"
	"github.com/konrad2002/tmate-server/dto"
	"github.com/konrad2002/tmate-server/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
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

	rg.POST("/user/login", uc.login)

	router.Use(auth.HandlerFunc(&uc.userService))

	router.GET("", uc.getAllUsers)
	router.GET("id/:id", uc.getUserById)
	router.GET("username/:username", uc.getUserByUsername)

	router.GET("me", uc.getUserForMe)

	router.POST("me/password", uc.changePasswordForMe)

	router.POST("password/:username", uc.changePasswordForUser)
	router.POST("", uc.createUser)

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
		fmt.Print(err.Error())
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

func (uc *UserController) getUserForMe(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if exists == false {
		err := errors.New("user subject not found")
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (uc *UserController) changePasswordForMe(c *gin.Context) {
	u, exists := c.Get("currentUser")
	if exists == false {
		err := errors.New("user subject not found")
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user := u.(dto.UserInfoDto)

	bodyAsByteArray, _ := ioutil.ReadAll(c.Request.Body)
	newPassword := string(bodyAsByteArray)

	if newPassword == "" {
		err := errors.New("empty password")
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	fmt.Printf("set %s's password to: %s\n", user.Username, newPassword)

	r, err := uc.userService.UpdatePassword(user.Username, newPassword, false)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func (uc *UserController) changePasswordForUser(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no username given"})
		return
	}

	bodyAsByteArray, _ := ioutil.ReadAll(c.Request.Body)
	newPassword := string(bodyAsByteArray)

	if newPassword == "" {
		err := errors.New("empty password")
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	fmt.Printf("set %s's password to: %s\n", username, newPassword)

	r, err := uc.userService.UpdatePassword(username, newPassword, true)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
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
