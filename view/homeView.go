package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../models"
	"time"
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

func (self HomeView) GetSettings(c *gin.Context, userUuid int) {
	c.HTML(http.StatusOK, "mySettings", gin.H{
		"Title": "Settings",
		"AccountUuid": userUuid,
		"Now": time.Now().Unix(),
	})
}
