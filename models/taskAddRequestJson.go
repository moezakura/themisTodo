package models

type TaskAddRequestJson struct {
	Name        string `json:"name" binding:"required"`
	Deadline    string `json:"deadline" binding:"required"`
	Assign      int    `json:"assign" binding:"required"`
	Description string `json:"description" binding:"required"`
	ProjectId   int    `json:"projectId"  binding:"required"`
}
