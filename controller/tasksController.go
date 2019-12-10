package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"themis.mox.si/themis/repository"
	"time"

	"github.com/gin-gonic/gin"
	"themis.mox.si/themis/models"
	"themis.mox.si/themis/module"
	"themis.mox.si/themis/utils"
	themisView "themis.mox.si/themis/view"
)

type TasksController struct {
	*BaseController
	projectRepo *repository.ProjectRepository
}

func NewTasksController(baseController *BaseController, projectRepo *repository.ProjectRepository) *TasksController {
	return &TasksController{BaseController: baseController, projectRepo: projectRepo}
}

func (t *TasksController) PostUpdate(c *gin.Context) {
	updateResult := &models.TaskUpdateResultJson{}
	createdTime, err := strconv.ParseInt(c.Param("createDate"), 10, 64)

	if err != nil {
		updateResult.Message = "invalid url"
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

	taskModule := module.NewTaskModule(t.DB)
	isErr, task := taskModule.Get(createdTime)

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
		if userTime.Unix() < now.Unix() && updateRequest.Deadline != task.Deadline {
			updateResult.Message = "deadline is not allowed past"
			themisView.TasksView{}.PostUpdate(c, updateResult)
			return
		}
	}

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
		isIn, err := t.projectRepo.IsIn(updateRequest.Assign, task.ProjectId)
		if err != nil {
			// TODO: エラーをログに出力する
			// TODO: エラーを返す
		}
		if !isIn {
			updateResult.Message = "invalid assign id"
			themisView.TasksView{}.PostUpdate(c, updateResult)
			return
		}
		task.Assign = updateRequest.Assign
	}

	uuid, _ := c.Get("uuid")
	taskModule.Update(createdTime, uuid.(int), task)

	updateResult.Success = true
	themisView.TasksView{}.PostUpdate(c, updateResult)
}

func (t *TasksController) PostBulkUpdate(c *gin.Context) {
	updateResult := &models.TaskUpdateResultJson{}

	var (
		taskStatusTarget   *models.TaskStatus = nil
		taskAssignTarget   int
		taskDeadlineTarget *time.Time = nil
	)
	projectsTarget := make([]int, 0)

	var updateRequest models.TaskBulkUpdateRequestJson
	c.ShouldBindJSON(&updateRequest)

	if len(updateRequest.BulkList) <= 0 {
		updateResult.Message = "no update list"
		themisView.TasksView{}.PostUpdate(c, updateResult)
		return
	}

	taskStatus := models.TaskStatus(updateRequest.Status)
	if taskStatus.String() != "OTHER" && taskStatus.String() != "Unknown" {
		taskStatusTarget = &taskStatus
	}

	taskModule := module.NewTaskModule(t.DB)

	isErr, tasks := taskModule.SearchCreateTimeList(updateRequest.BulkList)

	if isErr {
		updateResult.Message = "server error"
		themisView.TasksView{}.PostUpdate(c, updateResult)
		return
	}

	for _, task := range tasks {
		projectsTarget = append(projectsTarget, task.ProjectId)
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
		if userTime.Unix() > now.Unix() {
			updateResult.Message = "invalid deadline. deadline is not allow past."
			themisView.TasksView{}.PostUpdate(c, updateResult)
			return
		}

		taskDeadlineTarget = &userTime
	}

	if updateRequest.Assign > 0 {
		isIn, err := t.projectRepo.IsInBulk(updateRequest.Assign, projectsTarget)
		if err != nil {
			// TODO: エラーをログに出力する
			// TODO: エラーを返す
		}
		if !isIn {
			updateResult.Message = "invalid assign id"
			themisView.TasksView{}.PostUpdate(c, updateResult)
			return
		}
		taskAssignTarget = updateRequest.Assign
	}

	uuid, _ := c.Get("uuid")
	taskModule.UpdateAll(tasks, uuid.(int), taskStatusTarget, taskAssignTarget, taskDeadlineTarget)

	updateResult.Success = true
	themisView.TasksView{}.PostUpdate(c, updateResult)
}

func (t *TasksController) PostDelete(c *gin.Context) {
	deleteResult := &models.TaskDeleteResultJson{}
	createdTime, err := strconv.ParseInt(c.Param("createDate"), 10, 64)

	if err != nil {
		deleteResult.Message = "invalid task createdTime"
		themisView.TasksView{}.PostDelete(c, http.StatusBadRequest, deleteResult)
		return
	}

	userUuid := c.GetInt("uuid")

	taskModule := module.NewTaskModule(t.DB)
	isErr, task := taskModule.Get(createdTime)

	if isErr {
		deleteResult.Message = "invalid task createdTime"
		themisView.TasksView{}.PostDelete(c, http.StatusBadRequest, deleteResult)
		return
	}

	isIn, err := t.projectRepo.IsIn(userUuid, task.ProjectId)
	if err != nil {
		// TODO: エラーをログに出力する
		// TODO: エラーを返す
	}
	if !isIn {
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

func (t *TasksController) DeleteBulkDelete(c *gin.Context) {
	deleteResult := &models.TaskDeleteResultJson{}

	var deleteRequest models.TaskBulkDeleteRequestJson
	c.ShouldBindJSON(&deleteRequest)

	deleteList := make([]int64, 0)
	projectsTarget := make([]int, 0)

	if len(deleteRequest.BulkList) <= 0 {
		deleteResult.Message = "no delete list"
		themisView.TasksView{}.PostDelete(c, http.StatusBadRequest, deleteResult)
		return
	}

	for _, createdTimeString := range deleteRequest.BulkList {
		createdTime, err := strconv.ParseInt(createdTimeString, 10, 64)
		if err != nil {
			deleteResult.Message = "invalid task createdTime"
			themisView.TasksView{}.PostDelete(c, http.StatusBadRequest, deleteResult)
			return
		}
		deleteList = append(deleteList, createdTime)
	}

	userUuid := c.GetInt("uuid")

	taskModule := module.NewTaskModule(t.DB)
	isErr, tasks := taskModule.SearchCreateTimeList(deleteRequest.BulkList)
	if len(tasks) != len(deleteRequest.BulkList) {
		deleteResult.Message = "invalid task createdTime"
		themisView.TasksView{}.PostDelete(c, http.StatusBadRequest, deleteResult)
		return
	}
	if isErr {
		deleteResult.Message = "server error"
		themisView.TasksView{}.PostDelete(c, http.StatusInternalServerError, deleteResult)
		return
	}

	for _, task := range tasks {
		projectsTarget = append(projectsTarget, task.ProjectId)
	}

	isIn, err := t.projectRepo.IsInBulk(userUuid, projectsTarget)
	if err != nil {
		// TODO: エラーをログに出力する
		// TODO: エラーを返す
	}
	if !isIn {
		deleteResult.Message = "invalid token"
		themisView.TasksView{}.PostDelete(c, http.StatusBadRequest, deleteResult)
		return
	}

	var statusCode int
	if isErr := taskModule.DeleteAll(deleteList); isErr {
		deleteResult.Message = "delete failed"
		statusCode = http.StatusBadRequest
	} else {
		deleteResult.Success = true
		deleteResult.Message = ""
		statusCode = http.StatusOK
	}
	themisView.TasksView{}.PostDelete(c, statusCode, deleteResult)
}

func (t *TasksController) PostTaskCreate(c *gin.Context) {
	addResult := &models.TaskAddResultJson{}
	userUuid := c.GetInt("uuid")

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

	isIn, err := t.projectRepo.IsIn(addRequest.Assign, addRequest.ProjectId)
	if err != nil {
		// TODO: エラーをログに出力する
		// TODO: エラーを返す
	}
	if !isIn {
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

	taskModule := module.NewTaskModule(t.DB)
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

func (t *TasksController) GetView(c *gin.Context) {
	getResult := &models.TaskGetResultJson{}
	createdTime, err := strconv.ParseInt(c.Param("createDate"), 10, 64)
	if err != nil {
		getResult.Message = "invalid taskId"
		c.JSON(http.StatusOK, getResult)
		return
	}

	userUuid := c.GetInt("uuid")

	taskModule := module.NewTaskModule(t.DB)
	isErr, task := taskModule.Get(createdTime)
	if isErr {
		getResult.Message = "unknown taskId"
		themisView.TasksView{}.GetView(c, http.StatusBadRequest, getResult)
		return
	}

	isIn, err := t.projectRepo.IsIn(userUuid, task.ProjectId)
	if err != nil {
		// TODO: エラーをログに出力する
		// TODO: エラーを返す
	}
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

func (t *TasksController) GetSearch(c *gin.Context) {
	getResult := &models.TaskSearchResultJson{}
	searchRequest := models.TaskSearchRequest{}

	// taskIdがintであればsearchRequestにいれる
	taskIdTmp, err := strconv.ParseInt(c.Query("taskId"), 10, 64)
	if err == nil {
		searchRequest.TaskId = int(taskIdTmp)
	}

	// projectIdがintであればsearchRequestにいれる
	projectIdTmp, err := strconv.ParseInt(c.Query("projectId"), 10, 64)
	if err == nil {
		searchRequest.ProjectId = int(projectIdTmp)
	}

	// statusがintであればかつmodels.TaskStatusにパース成功すればsearchRequestにいれる
	statusTmp, err := strconv.ParseInt(c.Query("status"), 10, 64)
	if err == nil {
		taskStatus := models.TaskStatus(statusTmp)
		if taskStatus.String() != "OTHER" && taskStatus.String() != "Unknown" {
			searchRequest.Status = taskStatus
		}
	}

	// assignがintであればsearchRequestにいれる
	assignTmp, err := strconv.ParseInt(c.Query("assign"), 10, 64)
	if err == nil {
		searchRequest.AssignUserId = int(assignTmp)
	}

	// creatorがintであればsearchRequestにいれる
	createTmp, err := strconv.ParseInt(c.Query("creator"), 10, 64)
	if err == nil {
		searchRequest.CreateUserId = int(createTmp)
	}

	userUuid := c.GetInt("uuid")

	isIn, err := t.projectRepo.IsIn(userUuid, searchRequest.ProjectId)
	if err != nil {
		// TODO: エラーをログに出力する
		// TODO: エラーを返す
	}
	if !isIn {
		getResult.Message = "permission denied"
		themisView.TasksView{}.GetSearch(c, http.StatusForbidden, getResult)
		return
	}

	taskModule := module.NewTaskModule(t.DB)
	isErr, tasks := taskModule.Search(searchRequest)
	if isErr {
		getResult.Message = "unknown taskId"
		themisView.TasksView{}.GetSearch(c, http.StatusBadRequest, getResult)
		return
	}

	getResult.Tasks = make([]models.TaskOfJson, 0)
	getResult.Success = true
	for _, task := range tasks {
		taskTemp := utils.TaskConvert(&task)
		getResult.Tasks = append(getResult.Tasks, *models.NewTaskOfJson(*taskTemp))
	}

	themisView.TasksView{}.GetSearch(c, http.StatusOK, getResult)
}

func (t *TasksController) GetMy(c *gin.Context) {
	getResult := &models.TasksGetResultJson{}
	taskStatus, err := models.StringToTaskStatus(c.Query("status"))
	if err != nil {
		getResult.Message = "Unknown task status"
		themisView.TasksView{}.GetMy(c, http.StatusBadRequest, getResult)
		return
	}

	userUuid := c.GetInt("uuid")

	taskModule := module.NewTaskModule(t.DB)
	isErr, tasks := taskModule.GetTasksFromUser(userUuid, 20, taskStatus)
	if isErr {
		getResult.Message = "unknown taskId"
		themisView.TasksView{}.GetMy(c, http.StatusBadRequest, getResult)
		return
	}

	tasksTemp := utils.TasksConvert(tasks)
	getResult.Success = true
	getResult.Task = models.NewTasksOfJson(tasksTemp)

	themisView.TasksView{}.GetMy(c, http.StatusOK, getResult)
}

func (t *TasksController) GetHistoryList(c *gin.Context) {
	userUuid := c.GetInt("uuid")
	result := models.TaskHistoryJson{}
	result.Success = false
	createDate, err := strconv.ParseInt(c.Param("createDate"), 10, 64)

	if err != nil {
		result.Message = "invalid taskId"
		c.JSON(http.StatusBadRequest, result)
		return
	}

	tm := module.NewTaskModule(t.DB)

	isErr, task := tm.Get(createDate)
	if isErr {
		result.Message = "invalid taskId"
		c.JSON(http.StatusBadRequest, result)
		return
	}

	isIn, err := t.projectRepo.IsIn(userUuid, task.ProjectId)
	if err != nil {
		// TODO: エラーをログに出力する
		// TODO: エラーを返す
	}
	if !isIn {
		result.Message = "invalid taskId"
		c.JSON(http.StatusBadRequest, result)
		return
	}

	history, err := tm.GetHistoryList(task.CreateDate)
	if err != nil {
		result.Message = "server error"
		c.JSON(http.StatusServiceUnavailable, result)
		return
	}
	result.Success = true
	result.Message = ""
	result.Payload = make([]models.TaskHistoryOfJson, 0)
	for _, h := range history {
		deadline, err := time.Parse("2006-01-02T15:04:05-07:00", h.Task.Deadline)
		if err != nil {
			log.Printf("TasksModule.GetList Error: %+v\n", err)
		} else {
			deadlineFormatted := deadline.Format("2006-01-02")
			h.Task.Deadline = deadlineFormatted
		}
		h.Task = *utils.TaskHistoryItemConvert(&h.Task)
		result.Payload = append(result.Payload, h.ToJson())
	}
	c.JSON(http.StatusOK, result)
}

func (t *TasksController) PostApplyHistory(c *gin.Context) {
	userUuid := c.GetInt("uuid")
	_createDate := c.Param("createDate")

	res := models.TaskHistoryApplyResultJson{}
	var req models.TaskHistoryApplyRequestJson
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "invalid json"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	updateDate, err := strconv.ParseInt(req.UpdateDate, 10, 64)
	if err != nil {
		res.Message = "invalid updateDate"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	createDate, err := strconv.ParseInt(_createDate, 10, 64)
	if err != nil {
		res.Message = "invalid createDate"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	tm := module.NewTaskModule(t.DB)

	isErr, task := tm.Get(createDate)
	if isErr {
		res.Message = "not found createDate"
		c.JSON(http.StatusNotFound, res)
		return
	}

	isIn, err := t.projectRepo.IsIn(userUuid, task.ProjectId)
	if err != nil {
		// TODO: エラーをログに出力する
		// TODO: エラーを返す
	}
	if !isIn {
		res.Message = "not found createDate"
		c.JSON(http.StatusNotFound, res)
		return
	}

	if err := tm.ApplyHistory(createDate, updateDate); err != nil {
		res.Message = "history apply failed"
		c.JSON(http.StatusServiceUnavailable, res)
		return
	}

	res.Success = true
	res.Message = "ok"
	c.JSON(http.StatusOK, res)
}
