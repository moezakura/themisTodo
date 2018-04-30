package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../models"
)

type ProjectsView struct {
	*BaseView
}

func (self ProjectsView) GetAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "projectAdd", gin.H{
		"Title": "New Project",
	})
}

func (self ProjectsView) PostAdd(c *gin.Context, json *models.ProjectAddResultJson) {
	c.JSON(http.StatusOK, json)
}