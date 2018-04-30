package models

type ProjectAddResultJson struct {
	Success bool   `json:"success"`
	Id      int    `json:"id"`
	Message string `json:"message"`
}
