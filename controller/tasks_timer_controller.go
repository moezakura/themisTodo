package controller

import (
	"../models"
	"../module"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
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
