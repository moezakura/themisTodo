package models

import "strconv"

type TaskOfJson struct {
	TaskId      int        `json:"taskId"`
	ProjectId   int        `json:"projectId"`
	Name        string     `json:"name"`
	Creator     int        `json:"creator"`
	CreatorName string     `json:"creatorName"`
	Assign      int        `json:"assign"`
	AssignName  string     `json:"assignName"`
	Status      TaskStatus `json:"status"`
	Deadline    string     `json:"deadline"`
	LimitDate   int        `json:"limitDate"`
	DeadlineMD  string     `json:"deadlineMD"`
	Description string     `json:"description"`
	CreateDate  string     `json:"createDate"`
}

func NewTaskOfJson(task Task) *TaskOfJson {
	return &TaskOfJson{
		task.TaskId,
		task.ProjectId,
		task.Name,
		task.Creator,
		task.CreatorName,
		task.Assign,
		task.AssignName,
		task.Status,
		task.Deadline,
		task.LimitDate,
		task.DeadlineMD,
		task.Description,
		strconv.FormatInt(task.CreateDate, 10),
	}
}
