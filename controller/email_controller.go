package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konrad2002/tmate-server/dto"
	"github.com/konrad2002/tmate-server/service"
	"net/http"
)

type EmailController struct {
	emailService service.EmailService
	userService  service.UserService
}

func NewEmailController(emailService service.EmailService, userService service.UserService) EmailController {
	return EmailController{
		emailService: emailService,
		userService:  userService,
	}
}

func (emc *EmailController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/email/")

	router.GET("senders", emc.getEmailSenders)

	router.POST("send", emc.sendMail)

	router.OPTIONS("send", emc.ok)
}

func (emc *EmailController) ok(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (emc *EmailController) getEmailSenders(c *gin.Context) {
	fields, err := emc.emailService.GetEmailSenders()
	if err != nil {
		fmt.Print(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, fields)
}

func (emc *EmailController) sendMail(c *gin.Context) {
	var email dto.SendEmailDto
	if err := c.BindJSON(&email); err != nil {
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err := emc.emailService.SendEmailFromTemplate(email.Sender, email.Receivers, email.Subject, email.BodyTemplate)
	if err != nil {
		fmt.Print(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
