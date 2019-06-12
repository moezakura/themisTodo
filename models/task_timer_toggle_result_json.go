package models

type TaskTimerToggleResultJson struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Start   bool   `json:"start"`
}

func NewTaskTimerToggleResultJson(success bool) *TaskTimerToggleResultJson {
	return &TaskTimerToggleResultJson{Success: success}
}
