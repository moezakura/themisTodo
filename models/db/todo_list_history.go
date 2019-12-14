package db

import "time"

type TodoListHistory struct {
	Name        string    `gorm:"Column:name;type:VCHEAR(1024)"`
	Editor      int       `gorm:"Column:editor;type:INT(11)"`
	Status      int       `gorm:"Column:status;type:INT(11)"`
	Deadline    time.Time `gorm:"Column:deadline;type:DATE"`
	Description string    `gorm:"Column:description;type:TEXT"`
	CreateDate  int64     `gorm:"Column:createDate;type:BIGINT(20)"`
	UpdateDate  int64     `gorm:"Column:updateDate;type:BIGINT(20)"`
	Assign      int       `gorm:"Column:assign;type:INT(11)"`
}

func (t TodoListHistory) TableName() string {
	return "todo_list_history"
}
