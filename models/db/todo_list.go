package db

type TodoList struct {
	Id         int   `gorm:"Column:id;type:INT(11)"`
	Project    int   `gorm:"Column:project;type:INT(11)"`
	Creator    int   `gorm:"Column:creator;type:INT(11)"`
	CreateDate int64 `gorm:"Column:createDate;type:bigint(20);PRIMARY_KEY"`
	Adopted    int64 `gorm:"Column:adopted;type:bigint(20)"`
}

func (t TodoList) TableName() string {
	return "todo_list"
}
