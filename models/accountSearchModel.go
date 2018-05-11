package models

type AccountSearchModel struct {
	Name        string
	DisplayName string
	Uuid        int
	ProjectId   int
	Max         int
}

func NewAccountSearchModel() *AccountSearchModel {
	return &AccountSearchModel{
		"",
		"",
		-1,
		-1,
		10,
	}
}
