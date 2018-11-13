package models

type ProjectDeleteMemberResultJson struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Members []Account `json:"members"`
}
