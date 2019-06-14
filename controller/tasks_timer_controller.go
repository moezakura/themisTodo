package controller

import (
	"../models"
	"../module"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TaskTimerController struct {
	*BaseController
	watcher *module.TaskTimerWatcherModule
}

func NewTaskTimerController(baseController *BaseController, watcher *module.TaskTimerWatcherModule) *TaskTimerController {
	return &TaskTimerController{BaseController: baseController, watcher: watcher}
}


func (t *TaskTimerController) PatchToggle(c *gin.Context) {
	res := models.NewTaskTimerToggleResultJson(false)
	createdTime, err := strconv.ParseInt(c.Param("createDate"), 10, 64)

	if err != nil {
		res.Message = "invalid task createdTime"
		c.JSON(http.StatusServiceUnavailable, res)
		return
	}

	userUuid := c.GetInt("uuid")

	taskModule := module.NewTaskModule(t.DB)
	isErr, task := taskModule.Get(createdTime)

	if isErr {
		res.Message = "invalid task createdTime"
		c.JSON(http.StatusServiceUnavailable, res)
		return
	}

	projectModule := module.NewProjectsModule(t.DB)
	if isIn := projectModule.IsIn(userUuid, task.ProjectId); !isIn {
		res.Message = "invalid task createdTime"
		c.JSON(http.StatusServiceUnavailable, res)
		return
	}

	taskTimerModule := module.NewTasksTimerModule(t.DB, t.watcher)
	t.watcher.Job()
	isStart := false
	if isStart, err = taskTimerModule.TimerToggle(createdTime, userUuid); err != nil {
		log.Printf("%+v\n", err)
		res.Message = "Unknown Error"
		c.JSON(http.StatusServiceUnavailable, res)
		return
	}
	res.Success = true
	res.Start = isStart
	c.JSON(http.StatusOK, res)
}

func (t *TaskTimerController) GetView(c *gin.Context) {
	res := models.NewTaskTimerGetResult(false)
	createdTime, err := strconv.ParseInt(c.Param("createDate"), 10, 64)

	if err != nil {
		res.Message = "invalid task createdTime"
		c.JSON(http.StatusServiceUnavailable, res)
		return
	}

	userUuid := c.GetInt("uuid")

	taskModule := module.NewTaskModule(t.DB)
	isErr, task := taskModule.Get(createdTime)

	if isErr {
		res.Message = "invalid task createdTime"
		c.JSON(http.StatusServiceUnavailable, res)
		return
	}

	projectModule := module.NewProjectsModule(t.DB)
	if isIn := projectModule.IsIn(userUuid, task.ProjectId); !isIn {
		res.Message = "invalid task createdTime"
		c.JSON(http.StatusServiceUnavailable, res)
		return
	}

	taskTimerModule := module.NewTasksTimerModule(t.DB, t.watcher)
	histories, err := taskTimerModule.GetTaskTimerHistory(createdTime)

	res.TodayTime = 0
	res.TotalTime = 0
	res.LastStartTime = 0
	res.LastEndTime = 0
	res.Start = false

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local).Unix()

	for _, history := range histories {
		timeDiff := int(history.EndDateUnix - history.StartDateUnix)
		if history.EndDateUnix == 0 {
			res.Start = true
			timeDiff = int(time.Now().Unix() - history.StartDateUnix)
		}

		if history.StartDateUnix >= today && history.StartDateUnix <= todayEnd {
			res.TodayTime += timeDiff
		}
		res.TotalTime += timeDiff

		if history.StartDateUnix > res.LastStartTime {
			res.LastStartTime = history.StartDateUnix
			res.LastEndTime = history.EndDateUnix
		}
	}

	res.Success = true
	c.JSON(http.StatusOK, res)
}

func (t *TaskTimerController) GetMyList(c *gin.Context) {
	res := models.NewTaskTimerListResultJson(false)
	_projectId, err := strconv.ParseInt(c.Param("projectId"), 10, 64)
	startDateString := c.Query("start")
	endDateString := c.Query("end")

	if err != nil {
		res.Message = "invalid task createdTime"
		c.JSON(http.StatusServiceUnavailable, res)
		return
	}

	projectId := int(_projectId)
	userUuid := c.GetInt("uuid")

	projectModule := module.NewProjectsModule(t.DB)
	if isIn := projectModule.IsIn(userUuid, projectId); !isIn {
		res.Message = "invalid task createdTime"
		c.JSON(http.StatusServiceUnavailable, res)
		return
	}

	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Tokyo")
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)

	if startDateString != "" && endDateString != "" {
		todayStart, _ = time.Parse("2006-01-02 15:04:05", startDateString)
		todayEnd, _ = time.Parse("2006-01-02 15:04:05", endDateString)
	}

	taskTimerModule := module.NewTasksTimerModule(t.DB, t.watcher)
	histories, err := taskTimerModule.SearchTaskTimer([]int{projectId}, []int{userUuid}, &todayStart, &todayEnd)

	res.Success = true
	res.List = histories
	c.JSON(http.StatusOK, res)
}
