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
	"fmt"
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

	taskModule := module.NewTaskModule(self.DB, self.GormDB)
	projectModule := module.NewProjectsModule(self.DB)
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

	taskStatus := models.TaskStatus(updateRequest.Status)
	if taskStatus.String() != "OTHER" && taskStatus.String() != "Unknown" {
		task.Status = models.TaskStatus(updateRequest.Status)
	}

	if updateRequest.Assign > 0 {
		if !projectModule.IsIn(updateRequest.Assign, task.ProjectId) {
			updateResult.Message = "invalid assign id"
			themisView.TasksView{}.PostUpdate(c, updateResult)
			return
		}
		task.Assign = updateRequest.Assign
	}

	taskModule.Update(createdTime, task)

	updateResult.Success = true
	themisView.TasksView{}.PostUpdate(c, updateResult)
}

func (self TasksController) PostDelete(c *gin.Context) {
	deleteResult := &models.TaskDeleteResultJson{}
	createdTime, err := strconv.ParseInt(c.Param("createDate"), 10, 64)

	if err != nil {
		deleteResult.Message = "invalid task createdTime"
		themisView.TasksView{}.PostDelete(c, http.StatusBadRequest, deleteResult)
		return
	}

	loginModule := module.NewLoginModule(self.DB)
	isError, userUuid := loginModule.GetUserId(c, self.Session)

	if isError {
		deleteResult.Message = "invalid token"
		themisView.TasksView{}.PostDelete(c, http.StatusBadRequest, deleteResult)
		return
	}

	taskModule := module.NewTaskModule(self.DB, self.GormDB)
	isErr, task := taskModule.Get(createdTime)

	if isErr {
		deleteResult.Message = "invalid task createdTime"
		themisView.TasksView{}.PostDelete(c, http.StatusBadRequest, deleteResult)
		return
	}

	projectModule := module.NewProjectsModule(self.DB)
	if isIn := projectModule.IsIn(userUuid, task.ProjectId); !isIn {
		deleteResult.Message = "invalid task createdTime"
		themisView.TasksView{}.PostDelete(c, http.StatusBadRequest, deleteResult)
		return
	}

	var statusCode int
	if isErr := taskModule.Delete(createdTime); isErr {
		deleteResult.Message = "delete failed"
		statusCode = http.StatusBadRequest
	} else {
		deleteResult.Success = true
		deleteResult.Message = ""
		statusCode = http.StatusOK
	}
	themisView.TasksView{}.PostDelete(c, statusCode, deleteResult)
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

	projectModule := module.NewProjectsModule(self.DB)
	if !projectModule.IsIn(addRequest.Assign, addRequest.ProjectId) {
		addResult.Message = "invalid assign user id"
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

	taskModule := module.NewTaskModule(self.DB, self.GormDB)
	newTask := &models.Task{
		TaskId:      taskModule.GetLastId(addRequest.ProjectId) + 1,
		ProjectId:   addRequest.ProjectId,
		Name:        addRequest.Name,
		Creator:     userUuid,
		Assign:      addRequest.Assign,
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

	taskModule := module.NewTaskModule(self.DB, self.GormDB)
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

	taskTemp := utils.TaskConvert(task)
	getResult.Success = true
	getResult.Task = models.NewTaskOfJson(*taskTemp)

	themisView.TasksView{}.GetView(c, http.StatusOK, getResult)
}

func (self TasksController) GetSearch(c *gin.Context) {
	getResult := &models.TaskGetResultJson{}
	taskIdTmp, err := strconv.ParseInt(c.Query("taskId"), 10, 64)
	if err != nil {
		getResult.Message = "invalid taskId"
		c.JSON(http.StatusOK, getResult)
		return
	}
	taskId := int(taskIdTmp)
	projectIdTmp, err := strconv.ParseInt(c.Query("projectId"), 10, 64)
	if err != nil {
		getResult.Message = "invalid projectId"
		c.JSON(http.StatusOK, getResult)
		return
	}
	projectId := int(projectIdTmp)

	loginModule := module.NewLoginModule(self.DB)
	isError, userUuid := loginModule.GetUserId(c, self.Session)

	if isError {
		getResult.Message = "invalid token"
		themisView.TasksView{}.GetView(c, http.StatusBadRequest, getResult)
		return
	}

	taskModule := module.NewTaskModule(self.DB, self.GormDB)
	isErr, task := taskModule.GetFromTaskId(taskId, projectId)
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

	taskTemp := utils.TaskConvert(task)
	getResult.Success = true
	getResult.Task = models.NewTaskOfJson(*taskTemp)

	themisView.TasksView{}.GetView(c, http.StatusOK, getResult)
}

func (self TasksController) GetSearches(c *gin.Context) {
	searchesResult := &models.TaskSearchesResultJson{}

	loginModule := module.NewLoginModule(self.DB)
	isError, userUuid := loginModule.GetUserId(c, self.Session)

	if isError {
		searchesResult.Message = "invalid token"
		themisView.TasksView{}.GetSearches(c, http.StatusBadRequest, searchesResult)
		return
	}

	taskModule := module.NewTaskModule(self.DB, self.GormDB)

	queryName := c.Query("name")
	queryDescription := c.Query("description")

	statusText := c.Query("status")
	queryStatus, err := strconv.ParseInt(statusText, 10, 64)
	if err != nil {
		queryStatusTmp, parseErr := models.TaskStatusFromString(statusText)

		if parseErr != nil {
			searchesResult.Message = fmt.Sprintf("%s", parseErr)
			themisView.TasksView{}.GetSearches(c, http.StatusBadRequest, searchesResult)
			return
		}

		queryStatus = int64(queryStatusTmp)
	}

	queryProjectId, err := strconv.ParseInt(c.Query("projectId"), 10, 64)
	if err != nil {
		searchesResult.Message = "invalid project id"
		themisView.TasksView{}.GetSearches(c, http.StatusBadRequest, searchesResult)
		return
	}
	projectId := int(queryProjectId)

	projectModule := module.NewProjectsModule(self.DB)
	if !projectModule.IsIn(userUuid, projectId) {
		searchesResult.Message = "forbidden"
		themisView.TasksView{}.GetSearches(c, http.StatusBadRequest, searchesResult)
		return
	}

	tasks, err := taskModule.SearchTasks(queryName, queryDescription, int(queryStatus), projectId)

	searchesResult.Success = true
	searchesResult.Message = ""

	jsonTaskList := make([]models.TaskOfJson, 0)
	for _, value := range tasks {
		jsonTaskList = append(jsonTaskList, *models.NewTaskOfJson(value))
	}
	searchesResult.Task = jsonTaskList

	themisView.TasksView{}.GetSearches(c, http.StatusOK, searchesResult)
}
