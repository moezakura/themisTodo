package models

import "themis.mox.si/themis/models/db"

type ProjectInfoResultJson struct {
	db.Project
	Success bool   `json:"success"`
	Message string `json:"message"`
}
