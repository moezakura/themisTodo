package db

type Project struct {
	Uuid        int    `json:"uuid" gorm:"Column:uuid;PRIMARY_KEY;type:INT(11)"`
	Name        string `json:"name" gorm:"Column:name;VARCHAR(256)"`
	Description string `json:"description" gorm:"Column:description;VARCHAR(1024)"`
}
