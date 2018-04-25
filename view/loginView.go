package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginView struct {
}

func (self LoginView) GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}
