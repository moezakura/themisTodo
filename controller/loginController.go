package controller

import (
	"github.com/gin-gonic/gin"
	themisView "../view"
)

type LoginController struct {
}

func (self LoginController) GetLogin(c *gin.Context){
	themisView.LoginView{}.GetLogin(c)
}