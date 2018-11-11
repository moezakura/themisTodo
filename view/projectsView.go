package view

import (
	"../models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProjectsView struct {
	*BaseView
}


func (self ProjectsView) PostAdd(c *gin.Context, json *models.ProjectAddResultJson) {
	c.JSON(http.StatusOK, json)
}

func (self ProjectsView) PostTaskBoard(c *gin.Context, json *models.TaskAddResultJson) {
	c.JSON(http.StatusOK, json)
}

func (self ProjectsView) PostUpdate(c *gin.Context, json *models.ProjectAddResultJson) {
	c.JSON(http.StatusOK, json)
}

func (self ProjectsView) PostAddUser(c *gin.Context, json *models.ProjectAddResultJson) {
	c.JSON(http.StatusOK, json)
}

func (self ProjectsView) PostDeleteProject(c *gin.Context, statusCode int, json *models.ProjectDeleteResultJson) {
	c.JSON(statusCode, json)
}

func (self ProjectsView) GetInfo(c *gin.Context, statusCode int, json *models.ProjectInfoResultJson) {
	c.JSON(statusCode, json)
}

func (self ProjectsView) GetTasks(c *gin.Context, statusCode int, json *models.TasksGetResultJson) {
	c.JSON(statusCode, json)
}

func (self ProjectsView) GetMy(c *gin.Context, statusCode int, json *models.ProjectGetResultJson) {
	c.JSON(statusCode, json)
}