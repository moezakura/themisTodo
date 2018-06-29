package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../models"
)

type TasksView struct {
	*BaseView
}

func (self TasksView) PostUpdate(c *gin.Context, json *models.TaskUpdateResultJson) {
	c.JSON(http.StatusOK, json)
}

func (self TasksView) PostDelete(c *gin.Context, statusCode int, json *models.TaskDeleteResultJson) {
	c.JSON(http.StatusOK, json)
}

func (self TasksView) GetView(c *gin.Context, statusCode int, json *models.TaskGetResultJson) {
	c.JSON(statusCode, json)
}

func (self TasksView) GetMy(c *gin.Context, statusCode int, json *models.TasksGetResultJson) {
	c.JSON(statusCode, json)
}

func (self TasksView) GetProjects(c *gin.Context, statusCode int, json *models.ProjectGetResultJson) {
	c.JSON(statusCode, json)
}