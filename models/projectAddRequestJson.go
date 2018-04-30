package models

type ProjectAddRequestJson struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
