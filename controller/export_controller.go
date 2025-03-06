package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/konrad2002/tmate-server/auth"
	"github.com/konrad2002/tmate-server/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"time"
)

type ExportController struct {
	exportService service.ExportService
	userService   service.UserService
}

func NewExportController(exportService service.ExportService, userService service.UserService) ExportController {
	return ExportController{
		exportService: exportService,
		userService:   userService,
	}
}

func (ec *ExportController) RegisterRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/export/")

	router.Use(auth.HandlerFunc(&ec.userService))

	router.GET("excel/:queryId", ec.exportExcel)
}

func (ec *ExportController) exportExcel(c *gin.Context) {
	queryId, convErr := primitive.ObjectIDFromHex(c.Param("queryId"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	sortField := c.Query("sort_field")
	sortDirection, _ := strconv.Atoi(c.Query("sort_direction"))

	buf, err := ec.exportService.ExportFromQueryId(queryId, sortField, sortDirection)
	if err != nil {
		fmt.Print(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Set the correct headers for file download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=Mitglieder_%s.xlsx", time.Now().Format("2006-01-02_15-04-05")))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
