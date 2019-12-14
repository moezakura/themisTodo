package db

import (
	"time"
)

type UsersInProject struct {
	UserId     int       `gorm:"Column:user_id;type:INT(11)"`
	ProjectId  int       `gorm:"Column:project_id;type:INT(11)"`
	Enable     bool      `gorm:"Column:enable"`
	Expiration time.Time `gorm:"Column:expiration"`
}
