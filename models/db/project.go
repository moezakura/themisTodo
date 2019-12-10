package db

type Project struct {
	Uuid        int    `json:"uuid" gorm:"uuid;PRIMARY_KEY;type:INT(11)"`
	Name        string `json:"name" gorm:"name;VARCHAR(256)"`
	Description string `json:"description" gorm:"description;VARCHAR(1024)"`
}
