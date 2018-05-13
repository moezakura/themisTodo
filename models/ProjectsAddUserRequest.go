package models

type ProjectsAddUserRequest struct {
	Uuid int `json:"uuid" binding:"required"`
}
