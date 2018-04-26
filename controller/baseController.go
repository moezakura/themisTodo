package controller

import (
	"database/sql"
	"../module"
)

type BaseController struct {
	DB *sql.DB
	Session *module.SessionModule
}

func NewBaseController(db *sql.DB) *BaseController {
	return &BaseController{
		db,
		module.NewSessionModule(),
	}
}