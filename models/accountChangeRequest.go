package models

type AccountChangeRequestJson struct {
	Name            string `json:"name" binding:"required"`
	DisplayName     string `json:"displayName" binding:"required"`
	Uuid            int    `json:"uuid" binding:"required"`
	Password        string `json:"password" binding:"required"`
	CurrentPassword string `json:"currentPassword" binding:"required"`
}
