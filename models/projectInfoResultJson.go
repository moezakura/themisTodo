package models

type ProjectInfoResultJson struct {
	Project
	Success bool      `json:"success"`
	Message string    `json:"message"`
}
