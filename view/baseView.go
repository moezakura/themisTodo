package view

import (
	"github.com/gin-gonic/gin"
	"themis.mox.si/themis/module"
)

type BaseView struct {
	Router  *gin.Engine
	Session *module.SessionModule
}
