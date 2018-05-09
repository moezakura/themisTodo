package models

type Task struct {
	TaskId      int
	ProjectId   int
	Name        string
	Creator     int
	CreatorName string
	Status      TaskStatus
	Deadline    string
	LimitDate   int
	DeadlineMD  string
	Description string
	CreateDate  int64
}
