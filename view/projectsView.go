package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../models"
	"../utils"
)

type ProjectsView struct {
	*BaseView
}

func (self ProjectsView) GetAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "projectAdd", gin.H{
		"Title": "New Project",
	})
}

func (self ProjectsView) PostAdd(c *gin.Context, json *models.ProjectAddResultJson) {
	c.JSON(http.StatusOK, json)
}

func (self ProjectsView) GetTaskBoard(c *gin.Context, project *models.Project, taskList []models.Task, accounts []models.Account, creator *models.Account) {
	taskList = utils.TasksConvert(taskList)
	jsonTaskList := make([]models.TaskOfJson, 0)
	for _, value := range taskList  {
		jsonTaskList = append(jsonTaskList, *models.NewTaskOfJson(value))
	}

	c.HTML(http.StatusOK, "projectTaskBoard", gin.H{
		"Title":       project.Name,
		"Project":     project,
		"TaskList":    jsonTaskList,
		"AccountJson": accounts,
		"Creator":     creator,
	})
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
