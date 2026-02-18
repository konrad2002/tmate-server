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

type CourseController struct {
	courseService service.CourseService
	userService   service.UserService
}

func NewCourseController(courseService service.CourseService, userService service.UserService) CourseController {
	return CourseController{
		courseService: courseService,
		userService:   userService,
	}
}

func (cc *CourseController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/course/")

	router.Use(auth.HandlerFunc(&cc.userService))

	router.GET("", cc.getAllCourses)
	router.GET("id/:id", cc.getCourseById)
	router.GET("name/:name", cc.getCourseByName)
	router.POST("", cc.addCourse)
	router.PUT(":id", cc.updateCourse)
	router.DELETE(":id", cc.deleteCourse)
}

func (cc *CourseController) getAllCourses(c *gin.Context) {
	courses, err := cc.courseService.GetAll()
	if err != nil {
		fmt.Print(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, courses)
}

func (cc *CourseController) getCourseById(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	course, err := cc.courseService.GetById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "course not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, course)
}

func (cc *CourseController) getCourseByName(c *gin.Context) {
	name := c.Param("name")

	course, err := cc.courseService.GetByName(name)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "course not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, course)
}

func (cc *CourseController) addCourse(c *gin.Context) {
	u, _ := c.Get("currentUser")
	user := u.(dto.UserInfoDto)
	if !user.Permissions.TableStructureManagement {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "no table structure management permissions"})
		return
	}

	var course model.Course
	if err := c.BindJSON(&course); err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	r, err := cc.courseService.AddCourse(course)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, r)
}

func (cc *CourseController) updateCourse(c *gin.Context) {
	u, _ := c.Get("currentUser")
	user := u.(dto.UserInfoDto)
	if !user.Permissions.CourseManagement {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "no course management permissions"})
		return
	}

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	var course model.Course
	if err := c.BindJSON(&course); err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	course.Identifier = id

	r, err := cc.courseService.UpdateCourse(course)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func (cc *CourseController) deleteCourse(c *gin.Context) {
	u, _ := c.Get("currentUser")
	user := u.(dto.UserInfoDto)
	if !user.Permissions.CourseManagement {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "no course management permissions"})
		return
	}

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	err = cc.courseService.DeleteCourse(id)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusNoContent, nil)
}
