package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HomeView struct {
	*BaseView
}

func (self HomeView) GetHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home", gin.H{
		"Title": "Home",
	})
}
