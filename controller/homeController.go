package controller

import (
	"github.com/gin-gonic/gin"
	themisView "../view"
)

type HomeController struct {
	*BaseController
}

func (self HomeController) GetHome(c *gin.Context) {
	themisView.HomeView{self.BaseView}.GetHome(c)
}
