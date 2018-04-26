package view

import (
	"github.com/gin-gonic/gin"
	"../module"
)

type BaseView struct {
	Router  *gin.Engine
	Session *module.SessionModule
}
