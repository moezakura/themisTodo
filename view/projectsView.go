package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"themis.mox.si/themis/models"
	"themis.mox.si/themis/models/json"
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

func (self ProjectsView) GetInfo(c *gin.Context, statusCode int, json *json.ProjectInfoResultJson) {
	c.JSON(statusCode, json)
}

func (self ProjectsView) GetTasks(c *gin.Context, statusCode int, json *models.TasksGetResultJson) {
	c.JSON(statusCode, json)
}

func (self ProjectsView) GetMembers(c *gin.Context, statusCode int, json *models.ProjectMembersResultJson) {
	c.JSON(statusCode, json)
}

func (self ProjectsView) DeleteMember(c *gin.Context, statusCode int, json *models.ProjectDeleteMemberResultJson) {
	c.JSON(statusCode, json)
}

func (self ProjectsView) GetMy(c *gin.Context, statusCode int, json *json.ProjectGetResultJson) {
	c.JSON(statusCode, json)
}