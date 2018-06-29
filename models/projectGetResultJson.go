package models

type ProjectGetResultJson struct {
	Project []Project `json:"project"`
	Success bool      `json:"success"`
	Message string    `json:"message"`
}