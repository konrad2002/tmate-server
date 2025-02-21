package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func heatController() {
	router.GET("/heat", getHeats)
}

func getHeats(c *gin.Context) {
	heats, err := service.GetHeats()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, heats)
}
