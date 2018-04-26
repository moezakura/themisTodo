package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../models"
)

type LoginView struct {
}

func (self LoginView) GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (self LoginView) PostLogin(c *gin.Context, json *models.LoginResultJson) {
	c.JSON(http.StatusOK, json)
}