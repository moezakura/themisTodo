package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"log"
	"sync"
	"themis.mox.si/themis/models"
	"themis.mox.si/themis/models/db"
	"themis.mox.si/themis/module"
	"time"
)

type TaskRepository struct {
	db     *gorm.DB
	dbLock sync.Mutex
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		db:     db,
		dbLock: sync.Mutex{},
	}
}

func (t *TaskRepository) GetLastId(projectId int) (id int, err error) {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()

	todo := db.TodoList{}

	err = t.db.First(&todo, "project = ?", projectId).Error
	return todo.Id, err
}

func (t *TaskRepository) Add(task *models.Task) (taskResult *models.Task, err error) {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()

	tx := t.db.Begin()
	now := time.Now().UnixNano()
	loc, _ := time.LoadLocation("Asia/Tokyo")
	deadline, err := time.ParseInLocation("2006-01-02", task.Deadline, loc)

	todoHistory := db.TodoListHistory{
		Name:        task.Name,
		Editor:      task.Creator,
		Status:      int(task.Status),
		Deadline:    deadline,
		Description: task.Description,
		CreateDate:  now,
		UpdateDate:  now,
		Assign:      task.Assign,
	}

	err = tx.Create(&todoHistory).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	todo := db.TodoList{
		Id:         task.TaskId,
		Project:    task.ProjectId,
		Creator:    task.Creator,
		CreateDate: now,
		Adopted:    now,
	}

	err = tx.Create(&todo).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	task.CreateDate = now
	task.CreateDate = now

	return task, nil
}

func (t *TaskRepository) GetList(projectId int) (list []models.Task, err error) {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()
	list = []models.Task{}

	selects := t.getDefaultSelect()
	selectsText := t.makeSelect(selects)

	fullTodoList := make([]db.FullTodoList, 0)
	err = t.db.Table("todo_list").
		Select(selectsText).
		Joins("JOIN todo_list_history tlh on todo_list.adopted = tlh.updateDate").
		Joins("JOIN users u1 ON u1.uuid = todo_list.creator").
		Joins("JOIN users u2 ON u2.uuid = tlh.assign").
		Order("todo_list.id").
		Where("project = ?", projectId).
		Find(&fullTodoList, ).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		log.Printf("TasksModule.GetList Error: %+v\n", err)
		return nil, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return list, nil
	}

	for _, ftl := range fullTodoList {
		listOne := ftl.ToTask()

		list = append(list, listOne)
	}
	list, err = t.additionDoingStatus(list, []int{projectId})

	return list, err
}

func (t *TaskRepository) GetRow(createDate int64) (task *db.TodoList, err error) {
	selects := t.getDefaultSelect()
	selectsText := t.makeSelect(selects)

	fullTodoList := db.FullTodoList{}
	err = t.db.Table("todo_list").
		Select(selectsText).
		Joins("JOIN todo_list_history tlh on todo_list.adopted = tlh.updateDate").
		Joins("JOIN users u1 ON u1.uuid = todo_list.creator").
		Joins("JOIN users u2 ON u2.uuid = tlh.assign").
		Order("todo_list.id").
		First(&fullTodoList, "todo_list.createDate = ?", createDate).Error

	if err != nil {
		log.Printf("TasksModule.GetList Error: %+v\n", err)
		return nil, err
	}

	return &fullTodoList.TodoList, nil
}

func (t *TaskRepository) Get(createDate int64) (task *models.Task, err error) {
	selects := t.getDefaultSelect()
	selectsText := t.makeSelect(selects)

	fullTodoList := db.FullTodoList{}
	err = t.db.Table("todo_list").
		Select(selectsText).
		Joins("JOIN todo_list_history tlh on todo_list.adopted = tlh.updateDate").
		Joins("JOIN users u1 ON u1.uuid = todo_list.creator").
		Joins("JOIN users u2 ON u2.uuid = tlh.assign").
		Order("todo_list.id").
		First(&fullTodoList, "todo_list.createDate = ?", createDate).Error

	if err != nil {
		log.Printf("TasksModule.GetList Error: %+v\n", err)
		return nil, err
	}

	_task := fullTodoList.ToTask()
	task = &_task

	timer := module.NewTasksTimerModule(t.db.DB(), nil)

	isDoing, err := timer.GetTaskTimerStatus(createDate)
	if err != nil {
		log.Printf("TaskTimerModule.SearchTaskTimer Error: %+v\n", err)
		return nil, err
	}
	task.IsDoing = isDoing

	return task, nil
}

func (t *TaskRepository) GetBulk(projectId int, createDates []int64) (list []models.Task, err error) {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()

	list = []models.Task{}

	selects := t.getDefaultSelect()
	selectsText := t.makeSelect(selects)

	fullTodoList := make([]db.FullTodoList, 0)

	stmt := t.db.Table("todo_list").
		Select(selectsText).
		Joins("JOIN todo_list_history tlh on todo_list.adopted = tlh.updateDate").
		Joins("JOIN users u1 ON u1.uuid = todo_list.creator").
		Joins("JOIN users u2 ON u2.uuid = tlh.assign").
		Order("todo_list.id")

	if projectId > 0 {
		stmt = stmt.Where("todo_list.project = ?", projectId)
	}
	if len(createDates) > 0 {
		stmt = stmt.Where("todo_list.createDate IN (?)", createDates)
	}

	err = stmt.Find(&fullTodoList).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		log.Printf("TasksModule.GetList Error: %+v\n", err)
		return nil, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return list, nil
	}

	for _, ftl := range fullTodoList {
		listOne := ftl.ToTask()

		list = append(list, listOne)
	}

	list, err = t.additionDoingStatus(list, []int{projectId})

	return list, err
}

func (t *TaskRepository) Search(searchReq models.TaskSearchRequest) (list []models.Task, err error) {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()

	list = []models.Task{}

	selects := t.getDefaultSelect()
	selectsText := t.makeSelect(selects)

	fullTodoList := make([]db.FullTodoList, 0)

	stmt := t.db.Table("todo_list").
		Select(selectsText).
		Joins("JOIN todo_list_history tlh on todo_list.adopted = tlh.updateDate").
		Joins("JOIN users u1 ON u1.uuid = todo_list.creator").
		Joins("JOIN users u2 ON u2.uuid = tlh.assign").
		Order("todo_list.id")

	if searchReq.ProjectId > 0 {
		stmt = stmt.Where("todo_list.project = ?", searchReq.ProjectId)
	}

	if searchReq.TaskId > 0 {
		stmt = stmt.Where("todo_list.id = ?", searchReq.TaskId)
	}

	if searchReq.Status != models.TASK_STATUS_OTHER {
		stmt = stmt.Where("tlh.status = ?", int(searchReq.Status))
	}

	if searchReq.AssignUserId > 0 {
		stmt = stmt.Where("tlh.assign = ?", searchReq.AssignUserId)
	}

	if searchReq.CreateUserId > 0 {
		stmt = stmt.Where("todo_list.creator = ?", searchReq.CreateUserId)
	}
	err = stmt.Find(&fullTodoList).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		log.Printf("TasksModule.GetList Error: %+v\n", err)
		return nil, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return list, nil
	}
	projectIds := make([]int, 0)
	for _, ftl := range fullTodoList {
		listOne := ftl.ToTask()

		list = append(list, listOne)
		projectIds = append(projectIds, listOne.ProjectId)
	}

	list, err = t.additionDoingStatus(list, projectIds)

	return list, err
}

func (t *TaskRepository) SearchCreateTimeList(createDates []string) (list []models.Task, err error) {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()

	list = []models.Task{}

	selects := t.getDefaultSelect()
	selectsText := t.makeSelect(selects)

	fullTodoList := make([]db.FullTodoList, 0)
	err = t.db.Table("todo_list").
		Select(selectsText).
		Joins("JOIN todo_list_history tlh on todo_list.adopted = tlh.updateDate").
		Joins("JOIN users u1 ON u1.uuid = todo_list.creator").
		Joins("JOIN users u2 ON u2.uuid = tlh.assign").
		Order("todo_list.id").
		Find(&fullTodoList, "todo_list.createDate IN (?)", createDates).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		log.Printf("TasksModule.GetList Error: %+v\n", err)
		return nil, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return list, nil
	}

	projectIds := make([]int, 0)
	for _, ftl := range fullTodoList {
		listOne := ftl.ToTask()

		list = append(list, listOne)
		projectIds = append(projectIds, listOne.ProjectId)
	}
	list, err = t.additionDoingStatus(list, projectIds)

	return list, err
}

func (t *TaskRepository) Update(createDate int64, editor int, task *models.Task) (err error) {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()
	tx := t.db.Begin()

	now := time.Now().UnixNano()
	dadline, err := time.Parse("2006-01-02", task.Deadline)
	if err != nil {
		return err
	}

	history := db.TodoListHistory{
		Name:        task.Name,
		Editor:      editor,
		Status:      int(task.Status),
		Deadline:    dadline,
		Description: task.Description,
		CreateDate:  createDate,
		UpdateDate:  now,
		Assign:      task.Assign,
	}
	err = tx.Save(&history).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	todo, err := t.GetRow(createDate)
	if err != nil {
		tx.Rollback()
		return err
	}
	todo.Adopted = now
	err = tx.Save(todo).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (t *TaskRepository) UpdateAll(tasks []models.Task, editor int, status *models.TaskStatus, assign int, deadline *time.Time) (err error) {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()
	tx := t.db.Begin()

	createDates := make([]int64, 0)
	for _, t := range tasks {
		createDates = append(createDates, t.CreateDate)
	}

	todoList := make([]db.TodoList, 0)
	t.db.Find(&todoList, "createDate IN (?)", createDates)
	todoMap := make(map[int64]db.TodoList, 0)
	for _, t := range todoList {
		todoMap[t.CreateDate] = t
	}

	for _, task := range tasks {
		now := time.Now().UnixNano()
		updatedDeadline := deadline

		if status != nil {
			task.Status = *status
		}
		if assign > 0 {
			task.Assign = assign
		}
		if deadline == nil {
			_updatedDeadline, _ := time.Parse("2006-01-02", task.Deadline)
			updatedDeadline = &_updatedDeadline
		}

		history := db.TodoListHistory{
			Name:        task.Name,
			Editor:      editor,
			Status:      int(task.Status),
			Description: task.Description,
			CreateDate:  task.CreateDate,
			UpdateDate:  now,
			Assign:      task.Assign,
		}
		if updatedDeadline != nil {
			history.Deadline = *updatedDeadline
		}

		err = tx.Save(history).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		todo := todoMap[task.CreateDate]
		todo.Adopted = now

		err = tx.Save(todo).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if tx.Commit().Error != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (t *TaskRepository) Delete(createDate int64) (err error) {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()

	tx := t.db.Begin()

	err = tx.Exec("SET FOREIGN_KEY_CHECKS = 0;").Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Delete(&db.TodoList{}, "createDate = ?", createDate).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Delete(&db.TodoListHistory{}, "createDate = ?", createDate).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// TODO: todo timerの削除

	err = tx.Exec("SET FOREIGN_KEY_CHECKS = 1;").Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if tx.Commit().Error != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (t *TaskRepository) DeleteAll(createDates []int64) (err error) {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()

	tx := t.db.Begin()

	err = tx.Exec("SET FOREIGN_KEY_CHECKS = 0;").Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Delete(&db.TodoList{}, "createDate IN (?)", createDates).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Delete(&db.TodoListHistory{}, "createDate IN (?)", createDates).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Exec("SET FOREIGN_KEY_CHECKS = 1;").Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if tx.Commit() != nil {
		log.Printf("TasksModule.DeleteAll Error: (Transaction commit error) %+v\n", err)
	}

	return nil
}

func (t *TaskRepository) GetTasksFromUser(userUuid, limit int, status models.TaskStatus) (list []models.Task, err error) {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()
	list = []models.Task{}

	selects := t.getDefaultSelect()
	selectsText := t.makeSelect(selects)

	fullTodoList := make([]db.FullTodoList, 0)
	err = t.db.Table("todo_list").
		Select(selectsText).
		Joins("JOIN todo_list_history tlh on todo_list.adopted = tlh.updateDate").
		Joins("JOIN users u1 ON u1.uuid = todo_list.creator").
		Joins("JOIN users u2 ON u2.uuid = tlh.assign").
		Order("todo_list.id").
		Limit(limit).
		Find(&fullTodoList, "tlh.assign = ? AND tlh.status = ?", userUuid, int(status)).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		log.Printf("TasksModule.GetList Error: %+v\n", err)
		return nil, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return list, nil
	}

	projectIds := make([]int, 0)
	for _, ftl := range fullTodoList {
		listOne := ftl.ToTask()

		list = append(list, listOne)
		projectIds = append(projectIds, listOne.ProjectId)
	}
	list, err = t.additionDoingStatus(list, projectIds)

	return list, err
}

func (t *TaskRepository) GetHistoryList(createDate int64) (list []models.TaskHistory, err error) {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()

	selects := make(map[string][]string, 0)
	selects["todo_list_history"] = []string{"createDate", "name", "assign", "status", "deadline", "description", "updateDate", "editor"}
	selects["u1"] = []string{"displayName", "icon_path"}
	selects["u2"] = []string{"displayName", "icon_path"}
	selectsText := t.makeSelect(selects)

	fullTodoList := make([]db.FullTodoList, 0)
	err = t.db.Table("todo_list_history").
		Select(selectsText).
		Joins("JOIN users u1 ON u1.uuid = todo_list_history.editor").
		Joins("JOIN users u2 ON u2.uuid = todo_list_history.assign").
		Order("todo.id").
		Find(&fullTodoList, "todo_list_history.createDate = ?", createDate).Error
	if err != nil {
		return nil, err
	}

	list = make([]models.TaskHistory, 0)
	for _, h := range fullTodoList {
		tlh := h.TodoListHistory
		editor := h.User1
		assign := h.User2

		historyItem := models.TaskHistory{
			CreateDate: tlh.CreateDate,
			UpdateDate: tlh.UpdateDate,
			Task: models.TaskHistoryItem{
				Name:           tlh.Name,
				Editor:         tlh.Editor,
				EditorName:     editor.DisplayName,
				EditorIconPath: editor.IconPath,
				Status:         models.TaskStatus(tlh.Status),
				Assign:         tlh.Assign,
				AssignName:     assign.DisplayName,
				AssignIconPath: assign.IconPath,
				Deadline:       tlh.Deadline.Format("2006-01-02"),
				LimitDate:      0,
				DeadlineMD:     "",
				Description:    tlh.Description,
				CreateDate:     tlh.CreateDate,
				UpdateDate:     tlh.UpdateDate,
			},
		}
		list = append(list, historyItem)
	}

	return list, nil
}

func (t *TaskRepository) ApplyHistory(createDate, updateDate int64) error {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()

	c := 0
	err := t.db.Model(&db.TodoListHistory{}).
		Where("updateDate = ? AND createDate = ?", updateDate, createDate).
		Count(&c).Error

	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New("record not found")
	}

	task, err := t.GetRow(createDate)
	if err != nil {
		return err
	}
	tx := t.db.Begin()

	task.Adopted = updateDate
	err = tx.Update(task).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (t *TaskRepository) additionDoingStatus(tasks []models.Task, projectIds []int) (resultTasks []models.Task, err error) {
	timer := module.NewTasksTimerModule(t.db.DB(), nil)

	historys, err := timer.SearchTaskTimer(
		projectIds,
		[]int{},
		nil,
		nil,
		true,
	)
	if err != nil {
		return nil, err
	}
	for _, h := range historys {
		for k, t := range tasks {
			if t.CreateDate == h.CreateDate {
				tasks[k].IsDoing = true
				break
			}
		}
	}

	return tasks, nil
}

func (t *TaskRepository) getDefaultSelect() map[string][]string {
	selects := make(map[string][]string, 0)
	selects["todo_list"] = []string{"id", "project", "creator", "createDate", "adopted"}
	selects["tlh"] = []string{"createDate", "name", "assign", "status", "deadline", "description", "updateDate"}
	selects["u1"] = []string{"displayName", "icon_path"}
	selects["u2"] = []string{"displayName", "icon_path"}
	return selects
}

func (t *TaskRepository) makeSelect(keys map[string][]string) string {
	selectsText := ""
	for key, arr := range keys {
		for _, value := range arr {
			if len(selectsText) > 0 {
				selectsText += ","
			}
			selectsText += fmt.Sprintf("%s.%s", key, value)
		}
	}
	return selectsText
}
