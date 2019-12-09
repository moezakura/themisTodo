package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"themis.mox.si/themis/models"
)

type LoginView struct {
	*BaseView
}


func (self LoginView) PostLogin(c *gin.Context, json *models.LoginResultJson) {
	c.JSON(http.StatusOK, json)
}