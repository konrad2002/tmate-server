package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konrad2002/tmate-server/auth"
	"github.com/konrad2002/tmate-server/service"
	"net/http"
)

type AttestController struct {
	attestService service.AttestService
	userService   service.UserService
}

func NewAttestController(as service.AttestService, us service.UserService) AttestController {
	return AttestController{
		attestService: as,
		userService:   us,
	}
}

func (ac *AttestController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/attest/")

	router.Use(auth.HandlerFunc(&ac.userService))

	router.POST("exec", ac.executeRoutine)
}

func (ac *AttestController) executeRoutine(c *gin.Context) {
	err := ac.attestService.RunAttestRountine()

	if err != nil {
		fmt.Print(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
