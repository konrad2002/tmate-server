package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konrad2002/tmate-server/service"
	"net/http"
)

type MemberController struct {
	memberService service.MemberService
}

func NewMemberController(memberService service.MemberService) MemberController {
	return MemberController{
		memberService: memberService,
	}
}

func (mc *MemberController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/member/")
	router.GET("", mc.getAllMembers)
	router.GET("test", mc.getTest)
}

func (mc *MemberController) getTest(c *gin.Context) {
	test := mc.memberService.PrintTest()

	c.IndentedJSON(http.StatusOK, test)
}

func (mc *MemberController) getAllMembers(c *gin.Context) {
	members, err := mc.memberService.GetAll()
	if err != nil {
		fmt.Printf(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, members)
}
