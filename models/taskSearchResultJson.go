package models

type TaskSearchResultJson struct {
	Tasks   []TaskOfJson `json:"tasks"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
}
