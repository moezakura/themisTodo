package controller

import (
	"database/sql"
)

type BaseController struct {
	DB *sql.DB
}

func NewBaseController(db *sql.DB) *BaseController {
	return &BaseController{
		db,
	}
}