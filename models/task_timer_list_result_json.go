package models

type TaskTimerListResultJson struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	List    []TodoTimer `json:"list"`
}

func NewTaskTimerListResultJson(success bool) *TaskTimerListResultJson {
	return &TaskTimerListResultJson{Success: success}
}
