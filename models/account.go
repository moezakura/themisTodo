package models

type Account struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"displayName" binding:"required"`
	Uuid        int    `json:"uuid" binding:"required"`
	Password    string `json:"password" binding:"required"`
}
