package models

type AccountAddRequestJson struct {
	Name string `json:"name" binding:"required"`
}
