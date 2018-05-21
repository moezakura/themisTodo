package models

type Task struct {
	TaskId      int        `json:"taskId"`
	ProjectId   int        `json:"projectId"`
	Name        string     `json:"name"`
	Creator     int        `json:"creator"`
	CreatorName string     `json:"creatorName"`
	Status      TaskStatus `json:"status"`
	Deadline    string     `json:"deadline"`
	LimitDate   int        `json:"limitDate"`
	DeadlineMD  string     `json:"deadlineMD"`
	Description string     `json:"description"`
	CreateDate  int64      `json:"createDate"`
}
