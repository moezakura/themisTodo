package controller

import (
	"database/sql"
	"../module"
	"../view"
	"github.com/gin-gonic/gin"
)

type BaseController struct {
	DB *sql.DB
	Router *gin.Engine
	Session *module.SessionModule
	BaseView *view.BaseView
}

func NewBaseController(db *sql.DB, router *gin.Engine) *BaseController {
	session := module.NewSessionModule()
	return &BaseController{
		db,
		router,
		session,
		&view.BaseView{router, session},
	}
}