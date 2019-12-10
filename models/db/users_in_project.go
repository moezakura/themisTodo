package db

import (
	"github.com/jinzhu/gorm"
	"time"
)

type UsersInProject struct {
	gorm.DB
	UserId     int       `gorm:"user_id;type:INT(11)"`
	ProjectId  int       `gorm:"project_id;type:INT(11)"`
	Enable     bool      `gorm:"enable"`
	Expiration time.Time `gorm:"expiration"`
}
