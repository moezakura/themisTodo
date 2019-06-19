package models

type BaseApiResultJson struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewBaseApiResultJson(success bool) *BaseApiResultJson {
	return &BaseApiResultJson{Success: success}
}
