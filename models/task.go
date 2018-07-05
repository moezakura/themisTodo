package models

import (
	"github.com/jinzhu/gorm"
)

type Task struct {
	gorm.Model
	TaskId      int        `json:"taskId" gorm:"column:id"`
	ProjectId   int        `json:"projectId" gorm:"column:project"`
	Name        string     `json:"name"`
	Creator     int        `json:"creator"`
	CreatorName string     `json:"creatorName" gorm:"column:creatorName"`
	Assign      int        `json:"assign"`
	AssignName  string     `json:"assignName" gorm:"column:assignName"`
	Status      TaskStatus `json:"status"`
	Deadline    string     `json:"deadline"`
	LimitDate   int        `json:"limitDate"`
	DeadlineMD  string     `json:"deadlineMD"`
	Description string     `json:"description"`
	CreateDate  int64      `json:"createDate" gorm:"column:createDate;primary_key"`
}

func (Task) TableName() string {
	return "todo_list"
}
