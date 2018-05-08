package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../models"
)

type HomeView struct {
	*BaseView
}

func (self HomeView) GetHome(c *gin.Context, projects []models.Project) {
	c.HTML(http.StatusOK, "home", gin.H{
		"Title": "Home",
		"Projects": projects,
	})
}
