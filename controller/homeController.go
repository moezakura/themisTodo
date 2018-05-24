package controller

import (
	"github.com/gin-gonic/gin"
	themisView "../view"
	"../module"
	"../models"
	"net/http"
)

type HomeController struct {
	*BaseController
}

func (self HomeController) GetHome(c *gin.Context) {
	projectsModule := module.NewProjectsModule(self.DB)

	loginModule := module.NewLoginModule(self.DB)
	isError, userUuid := loginModule.GetUserId(c, self.Session)

	if isError {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	isError, projects := projectsModule.GetProjects(userUuid)
	if isError {
		c.String(http.StatusInternalServerError, "500 server error.")
		return
	}

	taskModule := module.NewTaskModule(self.DB)
	isErrorTodo, todoTaskList := taskModule.GetTasksFromUser(userUuid, 5, models.TASK_STATUS_TODO)
	isErrorDoing, doingTaskList := taskModule.GetTasksFromUser(userUuid, 5, models.TASK_STATUS_DOING)
	if isErrorTodo || isErrorDoing {
		c.String(http.StatusInternalServerError, "500 server error.")
		return
	}

	themisView.HomeView{self.BaseView}.GetHome(c, projects, todoTaskList, doingTaskList)
}

func (self HomeController) GetSettings(c *gin.Context) {

	loginModule := module.NewLoginModule(self.DB)
	isError, userUuid := loginModule.GetUserId(c, self.Session)

	if isError {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	themisView.HomeView{self.BaseView}.GetSettings(c, userUuid)
}
