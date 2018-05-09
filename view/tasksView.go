package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../models"
)

type TasksView struct {
	*BaseView
}

func (self TasksView) PostUpdate(c *gin.Context, json *models.TaskUpdateResultJson) {
	c.JSON(http.StatusOK, json)
}