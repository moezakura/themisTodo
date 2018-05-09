package controller

import (
	"github.com/gin-gonic/gin"
	themisView "../view"
)

type AccountController struct {
	*BaseController
}

func (self AccountController) GetAdd(c *gin.Context) {
	themisView.AccountView{self.BaseView}.GetAdd(c)
}
