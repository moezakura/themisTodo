package db

import "github.com/jinzhu/gorm"

type Project struct {
	gorm.Model
	Uuid        int    `json:"uuid" gorm:"uuid;type:INT(11)"`
	Name        string `json:"name" gorm:"name;VARCHAR(256)"`
	Description string `json:"description" gorm:"description;VARCHAR(1024)"`
}
