package models

import (
	"strconv"
)

type TaskHistoryJson struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Payload []TaskHistoryOfJson `json:"payload"`
}

type TaskHistoryOfJson struct {
	CreateDate string              `json:"create_date"`
	UpdateDate string              `json:"update_date"`
	Task       TaskHistoryItemJson `json:"task"`
}

type TaskHistory struct {
	CreateDate int64           `json:"create_date"`
	UpdateDate int64           `json:"update_date"`
	Task       TaskHistoryItem `json:"task"`
}

func (t TaskHistory) ToJson() TaskHistoryOfJson {
	return TaskHistoryOfJson{
		CreateDate: strconv.FormatInt(t.CreateDate, 10),
		UpdateDate: strconv.FormatInt(t.UpdateDate, 10),
		Task:       t.Task.ToJson(),
	}
}

type TaskHistoryItemJson struct {
	Name           string     `json:"name"`
	Editor         int        `json:"editor"`
	EditorName     string     `json:"editor_name"`
	EditorIconPath string     `json:"editor_icon_path"`
	Status         TaskStatus `json:"status"`
	Assign         int        `json:"assign"`
	AssignName     string     `json:"assign_name"`
	AssignIconPath string     `json:"assign_icon_path"`
	Deadline       string     `json:"deadline"`
	LimitDate      int        `json:"limit_date"`
	DeadlineMD     string     `json:"deadline_md"`
	Description    string     `json:"description"`
	CreateDate     string     `json:"create_date"`
	UpdateDate     string     `json:"update_date"`
}

type TaskHistoryItem struct {
	Name           string     `json:"name"`
	Editor         int        `json:"editor"`
	EditorName     string     `json:"editor_name"`
	EditorIconPath string     `json:"editor_icon_path"`
	Status         TaskStatus `json:"status"`
	Assign         int        `json:"assign"`
	AssignName     string     `json:"assign_name"`
	AssignIconPath string     `json:"assign_icon_path"`
	Deadline       string     `json:"deadline"`
	LimitDate      int        `json:"limit_date"`
	DeadlineMD     string     `json:"deadline_md"`
	Description    string     `json:"description"`
	CreateDate     int64      `json:"create_date"`
	UpdateDate     int64      `json:"update_date"`
}

func (t TaskHistoryItem) ToJson() TaskHistoryItemJson {
	return TaskHistoryItemJson{
		Name:           t.Name,
		Editor:         t.Editor,
		EditorName:     t.EditorName,
		EditorIconPath: t.EditorIconPath,
		Status:         t.Status,
		Assign:         t.Assign,
		AssignName:     t.AssignName,
		AssignIconPath: t.AssignIconPath,
		Deadline:       t.Deadline,
		LimitDate:      t.LimitDate,
		DeadlineMD:     t.DeadlineMD,
		Description:    t.Description,
		CreateDate:     strconv.FormatInt(t.CreateDate, 10),
		UpdateDate:     strconv.FormatInt(t.UpdateDate, 10),
	}
}
