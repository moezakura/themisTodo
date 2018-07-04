package module

import (
	"database/sql"
	"log"
	"../models"
	"time"
	"sync"
	"github.com/jinzhu/gorm"
	"fmt"
)

type TasksModule struct {
	db     *sql.DB
	gormDB *gorm.DB
	dbLock sync.Mutex
}

func NewTaskModule(db *sql.DB, gormDB *gorm.DB) *TasksModule {
	return &TasksModule{
		db,
		gormDB,
		sync.Mutex{},
	}
}

func (self *TasksModule) GetLastId(projectId int) int {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()
	rows, err := self.db.Query("SELECT `id` FROM `todo_list` WHERE `project` = ? ORDER BY `id` DESC LIMIT 0,1;", projectId)

	if err != nil {
		return 0
	}

	defer rows.Close()

	if !rows.Next() {
		return 0
	}

	var lastId int
	if err := rows.Scan(&lastId); err != nil {
		log.Printf("TasksModule.GetLastId Error: %+v\n", err)
		return 0
	}

	return lastId
}

func (self *TasksModule) Add(task *models.Task) *models.Task {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()

	stmt, err := self.db.Prepare("INSERT INTO `todo_list` (`id`, `project`, `name`, `creator`, `assign`, `status`, `deadline`, `description`, `createDate`) VALUE (?, ?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		log.Printf("TasksModule.Add Error: (1) %+v", err)
		return nil
	}

	defer stmt.Close()

	now := time.Now().UnixNano()
	_, err = stmt.Exec(task.TaskId, task.ProjectId, task.Name, task.Creator, task.Assign,
		task.Status, task.Deadline, task.Description, now)
	if err != nil {
		log.Printf("TasksModule.Add Error: (2) %+v", err)
		return nil
	}

	task.CreateDate = now

	return task
}

func (self *TasksModule) GetList(projectId int) (error bool, list []models.Task) {
	list = []models.Task{}

	rows, err := self.db.Query(`SELECT id, todo.name, creator, assign, status, deadline, description, createDate, u1.displayName, u2.displayName FROM todo_list todo
  INNER JOIN users u1 ON u1.uuid = todo.creator
  INNER JOIN users u2 ON u2.uuid = todo.assign
WHERE project = ? ORDER BY id ASC;`, projectId)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("TasksModule.GetList Error: %+v\n", err)
		return true, nil
	}

	if err == sql.ErrNoRows {
		return false, list
	}

	defer rows.Close()

	for rows.Next() {
		listOne := models.Task{}
		if err := rows.Scan(&listOne.TaskId, &listOne.Name, &listOne.Creator, &listOne.Assign,
			&listOne.Status, &listOne.Deadline, &listOne.Description,
			&listOne.CreateDate, &listOne.CreatorName, &listOne.AssignName); err != nil {
			log.Printf("TasksModule.GetList Error: %+v\n", err)
			return true, nil
		}
		list = append(list, listOne)
	}

	return false, list
}

func (self *TasksModule) Get(createDate int64) (isErr bool, task *models.Task) {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()
	rows, err := self.db.Query(`SELECT
  id,
  project,
  todo.name,
  creator,
  assign,
  status,
  deadline,
  description,
  createDate,
  u1.displayName,
  u2.displayName
FROM todo_list todo
  INNER JOIN users u1 ON u1.uuid = todo.creator
  INNER JOIN users u2 ON u2.uuid = todo.assign
WHERE createDate = ?;`, createDate)

	if err != nil {
		return true, nil
	}

	defer rows.Close()

	if !rows.Next() {
		return true, nil
	}

	returnTask := &models.Task{}
	if err := rows.Scan(&returnTask.TaskId, &returnTask.ProjectId, &returnTask.Name,
		&returnTask.Creator, &returnTask.Assign, &returnTask.Status, &returnTask.Deadline,
		&returnTask.Description, &returnTask.CreateDate, &returnTask.CreatorName, &returnTask.AssignName); err != nil {
		log.Printf("TasksModule.Get Error: %+v\n", err)
		return true, nil
	}

	return false, returnTask
}

func (self *TasksModule) GetFromTaskId(taskId int, projectId int) (isErr bool, task *models.Task) {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()
	rows, err := self.db.Query(`SELECT
  id,
  project,
  todo.name,
  creator,
  assign,
  status,
  deadline,
  description,
  createDate,
  u1.displayName,
  u2.displayName
FROM todo_list todo
  INNER JOIN users u1 ON u1.uuid = todo.creator
  INNER JOIN users u2 ON u2.uuid = todo.assign
WHERE todo.id = ? AND todo.project = ?;`, taskId, projectId)

	if err != nil {
		return true, nil
	}

	defer rows.Close()

	if !rows.Next() {
		return true, nil
	}

	returnTask := &models.Task{}
	if err := rows.Scan(&returnTask.TaskId, &returnTask.ProjectId, &returnTask.Name,
		&returnTask.Creator, &returnTask.Assign, &returnTask.Status, &returnTask.Deadline,
		&returnTask.Description, &returnTask.CreateDate, &returnTask.CreatorName, &returnTask.AssignName); err != nil {
		log.Printf("TasksModule.Get Error: %+v\n", err)
		return true, nil
	}

	return false, returnTask
}

func (self *TasksModule) Update(createDate int64, task *models.Task) (isErr bool) {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()
	_, err := self.db.Exec("UPDATE `todo_list` SET `name` = ?, `deadline` = ?, `description` = ?, `status` = ?, `assign` = ? WHERE `createDate` = ?;",
		task.Name, task.Deadline, task.Description, int(task.Status), task.Assign, createDate)

	if err != nil {
		log.Printf("TasksModule.Update Error: %+v\n", err)
		return true
	}

	return false
}

func (self *TasksModule) Delete(createDate int64) (isErr bool) {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()
	_, err := self.db.Exec("DELETE FROM `todo_list` WHERE `createDate` = ?;", createDate)

	if err != nil {
		log.Printf("TasksModule.Delete Error: %+v\n", err)
		return true
	}

	return false
}

func (self *TasksModule) GetTasksFromUser(userUuid, limit int, status models.TaskStatus) (isErr bool, list []models.Task) {
	list = []models.Task{}

	self.dbLock.Lock()
	defer self.dbLock.Unlock()
	rows, err := self.db.Query(`SELECT
  id,
  project,
  todo.name,
  creator,
  assign,
  status,
  deadline,
  description,
  createDate,
  u1.displayName,
  u2.displayName
FROM todo_list todo
  INNER JOIN users u1 ON u1.uuid = todo.creator
  INNER JOIN users u2 ON u2.uuid = todo.assign
WHERE assign = ? AND status = ?
LIMIT 0, ?;`, userUuid, status, limit)

	if err != nil {
		return true, nil
	}

	defer rows.Close()

	for rows.Next() {
		oneTask := models.Task{}
		if err := rows.Scan(&oneTask.TaskId, &oneTask.ProjectId, &oneTask.Name,
			&oneTask.Creator, &oneTask.Assign, &oneTask.Status, &oneTask.Deadline,
			&oneTask.Description, &oneTask.CreateDate, &oneTask.CreatorName, &oneTask.AssignName); err != nil {
			log.Printf("TasksModule.GetTasksFromUser Error: %+v\n", err)
			return true, nil
		}
		list = append(list, oneTask)
	}

	return false, list
}

func (self *TasksModule) SearchTasks(name, description string, status int, project int) (result []models.Task, err error) {
	sr := []models.Task{}
	db := self.gormDB.Table("todo_list")
	db = db.Select("id, project, todo_list.name,creator, assign, status, deadline," +
		"description, createDate,u1.displayName creatorName, u2.displayName assignName")

	searchORs := make([]string, 0)
	searchORsValue := make([]interface{}, 0)

	if len(name) > 0 {
		searchORs = append(searchORs, "todo_list.name LIKE ?")
		searchORsValue = append(searchORsValue, fmt.Sprintf("%%%s%%", name))
	}

	if len(description) > 0 {
		searchORs = append(searchORs, "description LIKE ?")
		searchORsValue = append(searchORsValue, fmt.Sprintf("%%%s%%", name))
	}

	orText := ""
	for k, v := range searchORs {
		t := v
		if k > 0 {
			t = " OR " + t
		}
		orText += t
	}
	db = db.Where(orText, searchORsValue...)

	statusAndProject := make([]string, 0)
	statusAndProjectValue := make([]interface{}, 0)
	if status > -1 {
		statusAndProject = append(statusAndProject, "status = ?")
		statusAndProjectValue = append(statusAndProjectValue, status)
	}

	if project > -1 {
		statusAndProject = append(statusAndProject, "project = ?")
		statusAndProjectValue = append(statusAndProjectValue, project)
	}

	for k, v := range statusAndProject {
		db = db.Where(v, statusAndProjectValue[k])
	}

	db = db.Joins("INNER JOIN users u1 ON u1.uuid = todo_list.creator")
	db = db.Joins("INNER JOIN users u2 ON u2.uuid = todo_list.assign")
	rows, err := db.Rows()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tsr := models.Task{}
		db.ScanRows(rows, &tsr)
		sr = append(sr, tsr)
	}

	return sr, nil
}
