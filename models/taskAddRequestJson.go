package models

type TaskAddRequestJson struct {
	Name        string `json:"name" binding:"required"`
	Deadline    string `json:"deadline" binding:"required"`
	Description string `json:"description" binding:"required"`
}
