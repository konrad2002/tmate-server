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
	"strconv"
)

type MemberController struct {
	memberService service.MemberService
	userService   service.UserService
}

func NewMemberController(memberService service.MemberService, userService service.UserService) MemberController {
	return MemberController{
		memberService: memberService,
		userService:   userService,
	}
}

func (mc *MemberController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/member/")

	router.Use(auth.HandlerFunc(&mc.userService))

	router.GET("", mc.getAllMembers)
	router.GET("id/:id", mc.getMemberById)
	router.GET("query/:queryId", mc.runMemberQuery)
	router.GET("families", mc.getFamilies)

	router.POST("", mc.addMember)
	router.POST("import", mc.importMembers)

	router.PUT("", mc.updateMember)

	router.OPTIONS("", mc.ok)

}

func (mc *MemberController) ok(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (mc *MemberController) getAllMembers(c *gin.Context) {
	members, err := mc.memberService.GetAll()
	if err != nil {
		fmt.Print(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, members)
}

func (mc *MemberController) getFamilies(c *gin.Context) {
	families, err := mc.memberService.GetFamilies()
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, families)
}

func (mc *MemberController) runMemberQuery(c *gin.Context) {
	queryId, convErr := primitive.ObjectIDFromHex(c.Param("queryId"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	sortField := c.Query("sort_field")
	sortDirection, _ := strconv.Atoi(c.Query("sort_direction"))

	members, fields, query, err := mc.memberService.GetAllByQueryId(queryId, sortField, sortDirection)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	result := dto.QueryResultDto{
		Members: *members,
		Fields:  *fields,
		Query:   *query,
	}

	c.IndentedJSON(http.StatusOK, result)
}

func (mc *MemberController) getMemberById(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	member, err := mc.memberService.GetById(id)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, member)
}

func (mc *MemberController) addMember(c *gin.Context) {
	var member model.Member
	if err := c.BindJSON(&member); err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var familyMemberId primitive.ObjectID
	if c.Query("family_member_id") != "" {
		var convErr error
		familyMemberId, convErr = primitive.ObjectIDFromHex(c.Query("family_member_id"))
		if convErr != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
			return
		}
	}

	r, err := mc.memberService.AddMember(member, familyMemberId)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func (mc *MemberController) importMembers(c *gin.Context) {
	var members []model.Member
	if err := c.BindJSON(&members); err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var errs []error
	success := 0
	for _, member := range members {
		_, err := mc.memberService.AddMember(member, primitive.ObjectID{0})
		if err == nil {
			success++
		} else {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			println(err.Error())
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"imported": success})
}

func (mc *MemberController) updateMember(c *gin.Context) {
	var member model.Member
	if err := c.BindJSON(&member); err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var familyMemberId primitive.ObjectID
	if c.Query("family_member_id") != "" {
		var convErr error
		familyMemberId, convErr = primitive.ObjectIDFromHex(c.Query("family_member_id"))
		if convErr != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
			return
		}
	}

	r, err := mc.memberService.UpdateMember(member, familyMemberId)
	if err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}
