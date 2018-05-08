package models

type Project struct {
	Uuid        int    `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
