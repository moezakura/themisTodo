package models

type TasksGetResultJson struct {
	Task    []TaskOfJson `json:"task"`
	Success bool          `json:"success"`
	Message string        `json:"message"`
}
