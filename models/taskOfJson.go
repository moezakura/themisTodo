package models

import "strconv"

type TaskOfJson struct {
	TaskId          int        `json:"taskId"`
	ProjectId       int        `json:"projectId"`
	Name            string     `json:"name"`
	Creator         int        `json:"creator"`
	CreatorName     string     `json:"creatorName"`
	CreatorIconPath string     `json:"creatorIconPath"`
	Assign          int        `json:"assign"`
	AssignName      string     `json:"assignName"`
	AssignIconPath  string     `json:"assignIconPath"`
	Status          TaskStatus `json:"status"`
	Deadline        string     `json:"deadline"`
	LimitDate       int        `json:"limitDate"`
	DeadlineMD      string     `json:"deadlineMD"`
	Description     string     `json:"description"`
	CreateDate      string     `json:"createDate"`
	Adopted         string     `json:"updateDate"`
}

func NewTaskOfJson(task Task) *TaskOfJson {
	return &TaskOfJson{
		task.TaskId,
		task.ProjectId,
		task.Name,
		task.Creator,
		task.CreatorName,
		task.CreatorIconPath,
		task.Assign,
		task.AssignName,
		task.AssignIconPath,
		task.Status,
		task.Deadline,
		task.LimitDate,
		task.DeadlineMD,
		task.Description,
		strconv.FormatInt(task.CreateDate, 10),
		strconv.FormatInt(task.Adopted, 10),
	}
}

func NewTasksOfJson(tasks []Task) []TaskOfJson {
	tasksOfJson := make([]TaskOfJson, 0)
	for _, task := range tasks {
		tasksOfJson = append(tasksOfJson, *NewTaskOfJson(task))
	}
	return tasksOfJson
}
