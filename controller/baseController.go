package controller

import (
	"database/sql"
	"../module"
	"../view"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type BaseController struct {
	DB *sql.DB
	GormDB *gorm.DB
	Router *gin.Engine
	Session *module.SessionModule
	BaseView *view.BaseView
}

func NewBaseController(db *sql.DB, gormDB *gorm.DB, router *gin.Engine) *BaseController {
	session := module.NewSessionModule()
	return &BaseController{
		db,
		gormDB,
		router,
		session,
		&view.BaseView{router, session},
	}
}