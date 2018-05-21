package controller

import (
	"github.com/gin-gonic/gin"
	themisView "../view"
	"../models"
	"../module"
	"../utils"
	"strconv"
	"strings"
	"time"
	"net/http"
)

type TasksController struct {
	*BaseController
}

func (self TasksController) PostUpdate(c *gin.Context) {
	updateResult := &models.TaskUpdateResultJson{}
	createdTime, err := strconv.ParseInt(c.Param("createDate"), 10, 64)

	if err != nil {
		updateResult.Message = "invalid url"
		themisView.TasksView{}.PostUpdate(c, updateResult)
		return
	}

	loginModule := module.NewLoginModule(self.DB)
	isError, _ := loginModule.GetUserId(c, self.Session)

	if isError {
		updateResult.Message = "invalid token"
		themisView.TasksView{}.PostUpdate(c, updateResult)
		return
	}

	var updateRequest models.TaskUpdateRequestJson
	c.ShouldBindJSON(&updateRequest)

	if len(updateRequest.Name) > 1024 {
		updateResult.Message = "maximum name length is 1024 characters"
		themisView.TasksView{}.PostUpdate(c, updateResult)
		return
	}

	if len(updateRequest.Description) > 10000 {
		updateResult.Message = "maximum description length is 10000 characters"
		themisView.TasksView{}.PostUpdate(c, updateResult)
		return
	}

	if len(updateRequest.Deadline) > 0 {
		timeSplits := strings.Split(updateRequest.Deadline, "-")
		if len(timeSplits) != 3 {
			updateResult.Message = "invalid deadline format(format is yyyy-mm-dd)"
			themisView.TasksView{}.PostUpdate(c, updateResult)
			return
		}

		errT, userTime := utils.ParseDate(updateRequest.Deadline)
		if errT {
			updateResult.Message = "invalid deadline format(date convert is failed)"
			themisView.TasksView{}.PostUpdate(c, updateResult)
			return
		}

		now := time.Now()
		if userTime.Unix() < now.Unix() {
			updateResult.Message = "deadline is not allowed past"
			themisView.TasksView{}.PostUpdate(c, updateResult)
			return
		}
	}

	taskModule := module.NewTaskModule(self.DB)
	isErr, task := taskModule.Get(createdTime)

	if isErr {
		updateResult.Message = "invalid task id"
		themisView.TasksView{}.PostUpdate(c, updateResult)
		return
	}

	if len(updateRequest.Name) > 0 {
		task.Name = updateRequest.Name
	}

	if len(updateRequest.Description) > 0 {
		task.Description = updateRequest.Description
	}

	if len(updateRequest.Deadline) > 0 {
		task.Deadline = updateRequest.Deadline
	}

	if updateRequest.Status > 0 || updateRequest.Status < 4 {
		task.Status = models.TaskStatus(updateRequest.Status)
	}

	taskModule.Update(createdTime, task)

	updateResult.Success = true
	themisView.TasksView{}.PostUpdate(c, updateResult)
}

func (self TasksController) PostTaskCreate(c *gin.Context) {
	addResult := &models.TaskAddResultJson{}
	loginModule := module.NewLoginModule(self.DB)
	isError, userUuid := loginModule.GetUserId(c, self.Session)

	if isError {
		addResult.Message = "invalid token"
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
		TaskId:      taskModule.GetLastId(addRequest.ProjectId) + 1,
		ProjectId:   addRequest.ProjectId,
		Name:        addRequest.Name,
		Creator:     userUuid,
		Status:      models.TASK_STATUS_TODO,
		Deadline:    addRequest.Deadline,
		Description: addRequest.Description,
		CreateDate:  0,
	}

	newTask = taskModule.Add(newTask)

	addResult.CreateDate = strconv.FormatInt(newTask.CreateDate, 10)
	addResult.Success = true
	themisView.ProjectsView{}.PostTaskBoard(c, addResult)
}

func (self TasksController) GetView(c *gin.Context) {
	getResult := &models.TaskGetResultJson{}
	createdTime, err := strconv.ParseInt(c.Param("createDate"), 10, 64)
	if err != nil {
		getResult.Message = "invalid taskId"
		c.JSON(http.StatusOK, getResult)
		return
	}

	loginModule := module.NewLoginModule(self.DB)
	isError, userUuid := loginModule.GetUserId(c, self.Session)

	if isError {
		getResult.Message = "invalid token"
		themisView.TasksView{}.GetView(c, http.StatusBadRequest, getResult)
		return
	}

	taskModule := module.NewTaskModule(self.DB)
	isErr, task := taskModule.Get(createdTime)
	if isErr {
		getResult.Message = "unknown taskId"
		themisView.TasksView{}.GetView(c, http.StatusBadRequest, getResult)
		return
	}

	projectsModule := module.NewProjectsModule(self.DB)
	isIn := projectsModule.IsIn(userUuid, task.ProjectId)
	if !isIn {
		getResult.Message = "permission denied"
		themisView.TasksView{}.GetView(c, http.StatusForbidden, getResult)
		return
	}

	getResult.Success = true
	getResult.Task = utils.TaskConvert(task)

	themisView.TasksView{}.GetView(c, http.StatusOK, getResult)
}
