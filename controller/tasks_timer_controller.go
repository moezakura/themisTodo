package controller

import (
	"../models"
	"../module"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"time"
)

type TaskTimerController struct {
	*BaseController
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

	taskTimerModule := module.NewTasksTimerModule(t.DB)
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

	taskTimerModule := module.NewTasksTimerModule(t.DB)
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
