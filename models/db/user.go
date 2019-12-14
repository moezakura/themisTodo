package db

type User struct {
	Uuid        int    `gorm:"Column:uuid;type:INT(11);PRIMARY_KEY"`
	DisplayName string `gorm:"Column:displayName;type:VARCHAR(256)"`
	Name        string `gorm:"Column:name;type:VARCHAR(128)"`
	IconPath    string `gorm:"Column:icon_path;type:VARCHAR(48)"`
	Password    string `gorm:"Column:password;type:VARCHAR(128)"`
}

func (t User) TableName() string {
	return "user"
}
