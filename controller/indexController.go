package controller

import (
	"github.com/gin-gonic/gin"
	themisView "../view"
)

type IndexController struct {
	*BaseController
}

func (self IndexController) GetIndex(c *gin.Context){
	themisView.IndexView{}.GetIndex(c)
}