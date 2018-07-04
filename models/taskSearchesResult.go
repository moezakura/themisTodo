package models

type TaskSearchesResultJson struct {
	Task    []TaskOfJson `json:"tasks"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
}
