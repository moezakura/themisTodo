package models

type TaskAddResultJson struct {
	Success    bool   `json:"success"`
	CreateDate int64  `json:"createDate"`
	Message    string `json:"message"`
}
