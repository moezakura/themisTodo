package models

type AccountSearchResultModel struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Uuid        int    `json:"uuid"`
}
