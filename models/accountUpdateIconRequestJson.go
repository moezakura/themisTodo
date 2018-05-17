package models

type AccountUpdateIconRequest struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	FileId string `json:"fileId"`
}
