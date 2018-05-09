package controller

import (
	"github.com/gin-gonic/gin"
	themisView "../view"
	"../utils"
	"../models"
	"../module"
	"net/http"
	"strconv"
	"time"
	"strings"
)

type ProjectsController struct {
	*BaseController
}

func (self ProjectsController) GetAdd(c *gin.Context) {
	themisView.ProjectsView{self.BaseView}.GetAdd(c)
}

func (self ProjectsController) PostAdd(c *gin.Context) {
	addResult := &models.ProjectAddResultJson{}
	loginModule := module.NewLoginModule(self.DB)
	isError, userUuid := loginModule.GetUserId(c, self.Session)

	if isError {
		addResult.Message = "invalid token"
		themisView.ProjectsView{}.PostAdd(c, addResult)
		return
	}

	var addRequest models.ProjectAddRequestJson
	c.ShouldBindJSON(&addRequest)

	if len(addRequest.Name) < 1 || len(addRequest.Description) < 1 {
		addResult.Message = "name and description is not allowed empty"
		themisView.ProjectsView{}.PostAdd(c, addResult)
		return
	}

	if len(addRequest.Name) > 256 {
		addResult.Message = "maximum name length is 256 characters"
		themisView.ProjectsView{}.PostAdd(c, addResult)
		return
	}

	if len(addRequest.Description) > 1024 {
		addResult.Message = "maximum description length is 1024 characters"
		themisView.ProjectsView{}.PostAdd(c, addResult)
		return
	}

	projectsModule := module.NewProjectsModule(self.DB)
	err2, id := projectsModule.Add(addRequest.Name, addRequest.Description)

	if err2 {
		addResult.Message = "server error"
		themisView.ProjectsView{}.PostAdd(c, addResult)
		return
	}
	projectsModule.AddUser(userUuid, id)

	addResult.Success = true
	addResult.Id = id
	themisView.ProjectsView{}.PostAdd(c, addResult)

	return
}

func (self ProjectsController) GetTaskBoard(c *gin.Context) {
	projectsModule := module.NewProjectsModule(self.DB)

	loginModule := module.NewLoginModule(self.DB)
	isError, _ := loginModule.GetUserId(c, self.Session)

	if isError {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	projectIdStr := c.Param("projectId")
	projectId64, err := strconv.ParseInt(projectIdStr, 10, 32)
	projectId := int(projectId64)
	if err != nil {
		c.String(http.StatusBadRequest, "400 Bad Request")
		return
	}

	isError, project := projectsModule.GetProject(projectId)
	taskModule := module.NewTaskModule(self.DB)

	iserr, taskList := taskModule.GetList(projectId)
	if iserr {
		c.String(http.StatusBadRequest, "400 Bad Request")
		return
	}

	themisView.ProjectsView{self.BaseView}.GetTaskBoard(c, project, taskList)
}

func (self ProjectsController) PostTaskBoard(c *gin.Context) {
	addResult := &models.TaskAddResultJson{}
	loginModule := module.NewLoginModule(self.DB)
	isError, userUuid := loginModule.GetUserId(c, self.Session)

	if isError {
		addResult.Message = "invalid token"
		themisView.ProjectsView{}.PostTaskBoard(c, addResult)
		return
	}

	projectIdStr := c.Param("projectId")
	projectId64, err := strconv.ParseInt(projectIdStr, 10, 32)
	projectId := int(projectId64)
	if err != nil {
		addResult.Message = "invalid project id"
		themisView.ProjectsView{}.PostTaskBoard(c, addResult)
		return
	}

	var addRequest models.TaskAddRequestJson
	c.ShouldBindJSON(&addRequest)

	if len(addRequest.Name) < 1 {
		addResult.Message = "name is not allowed empty"
		themisView.ProjectsView{}.PostTaskBoard(c, addResult)
		return
	}

	if len(addRequest.Name) > 1024 {
		addResult.Message = "maximum name length is 1024 characters"
		themisView.ProjectsView{}.PostTaskBoard(c, addResult)
		return
	}

	if len(addRequest.Description) > 10000 {
		addResult.Message = "maximum description length is 10000 characters"
		themisView.ProjectsView{}.PostTaskBoard(c, addResult)
		return
	}

	timeSplits := strings.Split(addRequest.Deadline, "-")
	if len(timeSplits) != 3 {
		addResult.Message = "invalid deadline format(format is yyyy-mm-dd)"
		themisView.ProjectsView{}.PostTaskBoard(c, addResult)
		return
	}

	errT, userTime := utils.ParseDate(addRequest.Deadline)
	if errT {
		addResult.Message = "invalid deadline format(date convert is failed)"
		themisView.ProjectsView{}.PostTaskBoard(c, addResult)
		return
	}

	now := time.Now()
	if userTime.Unix() < now.Unix() {
		addResult.Message = "deadline is not allowed past"
		themisView.ProjectsView{}.PostTaskBoard(c, addResult)
		return
	}

	taskModule := module.NewTaskModule(self.DB)
	newTask := &models.Task{
		TaskId:      taskModule.GetLastId(projectId) + 1,
		ProjectId:   projectId,
		Name:        addRequest.Name,
		Creator:     userUuid,
		Status:      models.TASK_STATUS_TODO,
		Deadline:    addRequest.Deadline,
		Description: addRequest.Description,
		CreateDate:  0,
	}

	newTask = taskModule.Add(newTask)

	addResult.CreateDate = newTask.CreateDate
	addResult.Success = true
	themisView.ProjectsView{}.PostTaskBoard(c, addResult)
}
