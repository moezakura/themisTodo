package models

type ProjectsDeleteUserRequest struct {
	Uuid int `json:"uuid" binding:"required"`
}
