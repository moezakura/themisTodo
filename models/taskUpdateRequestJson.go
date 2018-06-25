package models

type TaskUpdateRequestJson struct {
	Status      int    `json:"status" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Assign      int    `json:"assign" binding:"required"`
	Id          int    `json:"id" binding:"required"`
	Deadline    string `json:"deadline" binding:"required"`
	Description string `json:"description" binding:"required"`
}
