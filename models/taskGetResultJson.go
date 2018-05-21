package models

type TaskGetResultJson struct {
	Task    *Task  `json:"task"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}
