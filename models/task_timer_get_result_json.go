package models

type TaskTimerGetResult struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	Start         bool   `json:"start"`
	LastStartTime int64  `json:"last_start_time"`
	LastEndTime   int64  `json:"last_end_time"`
	TotalTime     int    `json:"total_time"`
	TodayTime     int    `json:"today_time"`
}

func NewTaskTimerGetResult(success bool) *TaskTimerGetResult {
	return &TaskTimerGetResult{Success: success}
}
