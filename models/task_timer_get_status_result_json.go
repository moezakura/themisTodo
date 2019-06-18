package models

type TaskTimerGetStatusResultJson struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Start   bool   `json:"start"`
}

func NewTaskTimerGetStatusResultJson(success bool) *TaskTimerGetStatusResultJson {
	return &TaskTimerGetStatusResultJson{Success: success}
}
