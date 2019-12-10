package models

import "themis.mox.si/themis/models/db"

type ProjectGetResultJson struct {
	Project []db.Project `json:"project"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
}
