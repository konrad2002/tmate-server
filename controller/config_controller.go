package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konrad2002/tmate-server/service"
	"net/http"
)

type ConfigController struct {
	configService service.ConfigService
}

func NewConfigController(configService service.ConfigService) ConfigController {
	return ConfigController{
		configService: configService,
	}
}

func (cc *ConfigController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/config/")

	router.GET("", cc.getConfig)
	router.GET("special_fields", cc.getSpecialFields)

	router.POST("init", cc.initConfig)
}

func (cc *ConfigController) getSpecialFields(c *gin.Context) {
	fields, err := cc.configService.GetSpecialFields()
	if err != nil {
		fmt.Printf(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, fields)
}

func (cc *ConfigController) getConfig(c *gin.Context) {
	config, err := cc.configService.GetConfig()
	if err != nil {
		fmt.Printf(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, config)
}

func (cc *ConfigController) initConfig(c *gin.Context) {
	err := cc.configService.InitConfig()

	if err != nil {
		fmt.Printf(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
