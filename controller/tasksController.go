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
