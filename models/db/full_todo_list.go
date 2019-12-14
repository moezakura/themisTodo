package db

import "themis.mox.si/themis/models"

type User1 User
type User2 User

type FullTodoList struct {
	TodoList        `gorm:"todo"`
	TodoListHistory `gorm:"tlh"`
	User1           `gorm:"u1"`
	User2           `gorm:"u2"`
}

func (f FullTodoList) ToTask() models.Task {
	todo := f.TodoList
	history := f.TodoListHistory
	creator := f.User1
	assign := f.User2

	return models.Task{
		TaskId:          todo.Id,
		ProjectId:       todo.Project,
		Name:            history.Name,
		Creator:         todo.Creator,
		CreatorName:     creator.DisplayName,
		CreatorIconPath: creator.IconPath,
		Assign:          history.Assign,
		AssignName:      assign.DisplayName,
		AssignIconPath:  assign.IconPath,
		Status:          models.TaskStatus(history.Status),
		Deadline:        history.Deadline.Format("2006-01-02"),
		LimitDate:       0,
		DeadlineMD:      "",
		Description:     history.Description,
		CreateDate:      history.CreateDate,
		Adopted:         todo.Adopted,
		IsDoing:         false,
	}
}
