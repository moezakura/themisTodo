package models

type ProjectMembersResultJson struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Members []Account `json:"members"`
}
