package models

type AccountListResult struct {
	Users   []Account `json:"users"`
	Success bool      `json:"success"`
	Message string    `json:"message"`
}
