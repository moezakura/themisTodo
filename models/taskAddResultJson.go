package models

type TaskAddResultJson struct {
	Success    bool   `json:"success"`
	CreateDate string `json:"createDate"`
	Message    string `json:"message"`
}
