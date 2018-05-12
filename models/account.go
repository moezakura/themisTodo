package models

type Account struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Uuid        int    `json:"uuid"`
	Password    string `json:"password"`
}
