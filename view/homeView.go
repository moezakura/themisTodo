package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../models"
	"../utils"
	"time"
)

type HomeView struct {
	*BaseView
}

func (self HomeView) GetHome(c *gin.Context, projects []models.Project, todoTaskList, doingTaskList []models.Task) {
	todoTaskList = utils.TasksConvert(todoTaskList)
	doingTaskList = utils.TasksConvert(doingTaskList)

	jsonTodoTaskList := make([]models.TaskOfJson, 0)
	for _, value := range todoTaskList  {
		jsonTodoTaskList = append(jsonTodoTaskList, *models.NewTaskOfJson(value))
	}

	jsonDoingTaskList := make([]models.TaskOfJson, 0)
	for _, value := range doingTaskList  {
		jsonDoingTaskList = append(jsonDoingTaskList, *models.NewTaskOfJson(value))
	}

	c.HTML(http.StatusOK, "home", gin.H{
		"Title":    "Home",
		"Projects": projects,
		"Todo":     jsonTodoTaskList,
		"Doing":    jsonDoingTaskList,
	})
}

func (self HomeView) GetSettings(c *gin.Context, userUuid int) {
	c.HTML(http.StatusOK, "mySettings", gin.H{
		"Title":       "Settings",
		"AccountUuid": userUuid,
		"Now":         time.Now().Unix(),
	})
}
