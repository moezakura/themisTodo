package models

type AccountProfileResultJson struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	User    *Account `json:"user"`
}
