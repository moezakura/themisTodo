package module

import (
	"../models"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"strings"
	"sync"
	"time"
)

type TasksModule struct {
	db     *sql.DB
	dbLock sync.Mutex
}

func NewTaskModule(db *sql.DB) *TasksModule {
	return &TasksModule{
		db,
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

	tx, err := self.db.Begin()
	if err != nil {
		log.Printf("TasksModule.Add Error: (txError) %+v", err)
		return nil
	}
	// insert list
	listStmt, err := tx.Prepare("INSERT INTO `todo_list` (`id`, `project`,`creator`,`createDate`) VALUE (?, ?, ?, ?)")
	if err != nil {
		log.Printf("TasksModule.Add Error: (Statement error) %+v\n", err)
		err := tx.Rollback()
		if err != nil {
			log.Printf("TasksModule.Add Error: (Transaction rolback error) %+v\n", err)
		}
		return nil
	}

	defer func() {
		err := listStmt.Close()
		if err != nil {
			log.Printf("TasksModule.Add Error: (Statement close error) %+v\n", err)
		}
	}()

	now := time.Now().UnixNano()
	_, err = listStmt.Exec(task.TaskId, task.ProjectId, task.Creator, now)
	if err != nil {
		log.Printf("TasksModule.Add Error: (Statement exec error) %+v", err)
		err := tx.Rollback()
		if err != nil {
			log.Printf("TasksModule.Add Error: (Transaction rolback error) %+v\n", err)
		}
		return nil
	}

	// insert history
	historyStmt, err := tx.Prepare("INSERT INTO `todo_list_history` (name, editor, status, deadline, description, createDate, updateDate, assign) VALUE (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("TasksModule.Add Error: (Statement error) %+v\n", err)
		err := tx.Rollback()
		if err != nil {
			log.Printf("TasksModule.Add Error: (Transaction rolback error) %+v\n", err)
		}
		return nil
	}
	defer func() {
		err := historyStmt.Close()
		if err != nil {
			log.Printf("TasksModule.Add Error: (Statement close error) %+v\n", err)
		}
	}()
	_, err = historyStmt.Exec(task.Name, task.Creator, task.Status, task.Deadline, task.Description, now, now, task.Assign)
	if err != nil {
		log.Printf("TasksModule.Add Error: (Statement exec error) %+v", err)
		err := tx.Rollback()
		if err != nil {
			log.Printf("TasksModule.Add Error: (Transaction rolback error) %+v\n", err)
		}
		return nil
	}

	// update list (connect history)
	listConnectStmt, err := tx.Prepare("UPDATE `todo_list` SET `adopted` = ? WHERE createDate = ?;")
	if err != nil {
		log.Printf("TasksModule.Add Error: (Statement error) %+v\n", err)
		err := tx.Rollback()
		if err != nil {
			log.Printf("TasksModule.Add Error: (Transaction rolback error) %+v\n", err)
		}
		return nil
	}

	defer func() {
		err := listConnectStmt.Close()
		if err != nil {
			log.Printf("TasksModule.Add Error: (Statement close error) %+v\n", err)
		}
	}()

	_, err = listConnectStmt.Exec(now, now)
	if err != nil {
		log.Printf("TasksModule.Add Error: (Statement exec error) %+v", err)
		err := tx.Rollback()
		if err != nil {
			log.Printf("TasksModule.Add Error: (Transaction rolback error) %+v\n", err)
		}
		return nil
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("TasksModule.Add Error: (Transaction commit error) %+v\n", err)
	}

	task.CreateDate = now

	return task
}

func (self *TasksModule) GetList(projectId int) (error bool, list []models.Task) {
	list = []models.Task{}

	rows, err := self.db.Query(`SELECT id, tlh.name, creator, tlh.assign, tlh.status, tlh.deadline, tlh.description, todo.createDate, todo.adopted, u1.displayName, u1.icon_path, u2.displayName, u2.icon_path FROM todo_list todo
  INNER JOIN todo_list_history tlh on todo.adopted = tlh.updateDate
  INNER JOIN users u1 ON u1.uuid = todo.creator
  INNER JOIN users u2 ON u2.uuid = tlh.assign
WHERE project = ? ORDER BY id ASC;`, projectId)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("TasksModule.GetList Error: %+v\n", err)
		return true, nil
	}

	if err == sql.ErrNoRows {
		return false, list
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("TasksModule.GetList Error: (rows close) %+v\n", err)
		}
	}()

	for rows.Next() {
		listOne := models.Task{}
		if err := rows.Scan(&listOne.TaskId, &listOne.Name, &listOne.Creator, &listOne.Assign,
			&listOne.Status, &listOne.Deadline, &listOne.Description, &listOne.CreateDate,
			&listOne.Adopted, &listOne.CreatorName, &listOne.CreatorIconPath, &listOne.AssignName,
			&listOne.AssignIconPath); err != nil {
			log.Printf("TasksModule.GetList Error: %+v\n", err)
			return true, nil
		}
		deadline, err := time.Parse("2006-01-02T15:04:05-07:00", listOne.Deadline)
		if err != nil {
			log.Printf("TasksModule.GetList Error: %+v\n", err)
		} else {
			deadlineFormatted := deadline.Format("2006-01-02")
			listOne.Deadline = deadlineFormatted
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
  tlh.name,
  creator,
  tlh.assign,
  tlh.status,
  tlh.deadline,
  tlh.description,
  todo.createDate,
  todo.adopted,
  u1.displayName,
  u1.icon_path,
  u2.displayName,
  u2.icon_path
FROM todo_list todo
  INNER JOIN todo_list_history tlh on todo.adopted = tlh.updateDate
  INNER JOIN users u1 ON u1.uuid = todo.creator
  INNER JOIN users u2 ON u2.uuid = tlh.assign
WHERE todo.createDate = ?;`, createDate)

	if err != nil {
		return true, nil
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("TasksModule.Get Error: (rows close) %+v\n", err)
		}
	}()

	if !rows.Next() {
		return true, nil
	}

	returnTask := &models.Task{}
	if err := rows.Scan(&returnTask.TaskId, &returnTask.ProjectId, &returnTask.Name,
		&returnTask.Creator, &returnTask.Assign, &returnTask.Status, &returnTask.Deadline,
		&returnTask.Description, &returnTask.CreateDate, &returnTask.Adopted, &returnTask.CreatorName,
		&returnTask.CreatorIconPath, &returnTask.AssignName, &returnTask.AssignIconPath); err != nil {
		log.Printf("TasksModule.Get Error: %+v\n", err)
		return true, nil
	} else {
		deadline, err := time.Parse("2006-01-02T15:04:05-07:00", returnTask.Deadline)
		if err != nil {
			log.Printf("TasksModule.GetList Error: %+v\n", err)
		} else {
			deadlineFormatted := deadline.Format("2006-01-02")
			returnTask.Deadline = deadlineFormatted
		}
	}

	return false, returnTask
}

func (self *TasksModule) GetBulk(projectId int, createDates []int64) (tasks []*models.Task, err error) {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()

	args := make([]interface{}, 0)
	args = append(args, projectId)
	for _, item := range createDates {
		args = append(args, item)
	}

	if len(createDates) == 0 {
		return make([]*models.Task, 0), nil
	}

	queryString := `SELECT
  id,
  project,
  tlh.name,
  creator,
  tlh.assign,
  tlh.status,
  tlh.deadline,
  tlh.description,
  todo.createDate,
  todo.adopted,
  u1.displayName,
  u1.icon_path,
  u2.displayName,
  u2.icon_path
FROM todo_list todo
  INNER JOIN todo_list_history tlh on todo.adopted = tlh.updateDate
  INNER JOIN users u1 ON u1.uuid = todo.creator
  INNER JOIN users u2 ON u2.uuid = tlh.assign
WHERE todo.project = ? AND todo.createDate IN (?` + strings.Repeat(", ?", len(createDates)-1) + `);`
	rows, err := self.db.Query(queryString, args...)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("TasksModule.GetBulk Error: (rows close) %+v\n", err)
		}
	}()

	tasks = make([]*models.Task, 0)
	for rows.Next() {
		returnTask := &models.Task{}
		if err := rows.Scan(&returnTask.TaskId, &returnTask.ProjectId, &returnTask.Name,
			&returnTask.Creator, &returnTask.Assign, &returnTask.Status, &returnTask.Deadline,
			&returnTask.Description, &returnTask.CreateDate, &returnTask.Adopted, &returnTask.CreatorName,
			&returnTask.CreatorIconPath, &returnTask.AssignName, &returnTask.AssignIconPath); err != nil {
			log.Printf("TasksModule.GetBulk Error: %+v\n", err)
			continue
		}

		deadline, err := time.Parse("2006-01-02T15:04:05-07:00", returnTask.Deadline)
		if err != nil {
			log.Printf("TasksModule.GetList Error: %+v\n", err)
		} else {
			deadlineFormatted := deadline.Format("2006-01-02")
			returnTask.Deadline = deadlineFormatted
		}

		tasks = append(tasks, returnTask)
	}

	return tasks, nil
}

func (self *TasksModule) Search(searchReq models.TaskSearchRequest) (isErr bool, task []models.Task) {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()

	joinString, whereString, queryArray := searchReq.ToSqlQueryAndArgs("tlh", "todo")
	joinQuery := `INNER JOIN todo_list_history tlh on todo.adopted = tlh.updateDate`
	if len(joinString) > 0 {
		joinQuery += ` AND ` + joinString
	}
	query := `SELECT
  id,
  project,
  tlh.name,
  creator,
  tlh.assign,
  tlh.status,
  tlh.deadline,
  tlh.description,
  todo.createDate,
  todo.adopted,
  u1.displayName,
  u1.icon_path,
  u2.displayName,
  u2.icon_path,
  tlh.updateDate
FROM todo_list todo
  ` + joinQuery + `
  INNER JOIN users u1 ON u1.uuid = todo.creator
  INNER JOIN users u2 ON u2.uuid = tlh.assign
WHERE ` + whereString + ";"

	rows, err := self.db.Query(query, queryArray...)

	if err != nil {
		log.Printf("TasksModule.Search Error:(query error) %+v\n", err)
		return true, nil
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("TasksModule.Search Error: (rows close) %+v\n", err)
		}
	}()

	returnTask := make([]models.Task, 0)

	for rows.Next() {
		task := models.Task{}
		var updateDate int64
		if err := rows.Scan(&task.TaskId, &task.ProjectId, &task.Name, &task.Creator, &task.Assign,
			&task.Status, &task.Deadline, &task.Description, &task.CreateDate, &task.Adopted, &task.CreatorName,
			&task.CreatorIconPath, &task.AssignName, &task.AssignIconPath, &updateDate); err != nil {
			log.Printf("TasksModule.Search Error:(scan error) %+v\n", err)
			return true, nil
		}
		if updateDate == 0 {
			fmt.Println("SKIP!")
			continue
		}

		deadline, err := time.Parse("2006-01-02T15:04:05-07:00", task.Deadline)
		if err != nil {
			log.Printf("TasksModule.GetList Error: %+v\n", err)
		} else {
			deadlineFormatted := deadline.Format("2006-01-02")
			task.Deadline = deadlineFormatted
		}
		returnTask = append(returnTask, task)
	}

	return false, returnTask
}

func (self *TasksModule) SearchCreateTimeList(searchReq []string) (isErr bool, task []models.Task) {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()

	queryString := ""
	queryArray := make([]interface{}, 0)
	for _, value := range searchReq {
		if len(queryString) > 0 {
			queryString += " OR "
		}
		queryString += " todo.createDate = ? "
		queryArray = append(queryArray, value)
	}

	rows, err := self.db.Query(`SELECT
  id,
  project,
  tlh.name,
  creator,
  tlh.assign,
  tlh.status,
  tlh.deadline,
  tlh.description,
  todo.createDate,
  todo.adopted,
  u1.displayName,
  u1.icon_path,
  u2.displayName,
  u2.icon_path
FROM todo_list todo
  INNER JOIN todo_list_history tlh on todo.adopted = tlh.updateDate
  INNER JOIN users u1 ON u1.uuid = todo.creator
  INNER JOIN users u2 ON u2.uuid = tlh.assign
WHERE `+queryString+";", queryArray...)

	if err != nil {
		log.Printf("TasksModule.SearchCreateTimeList Error:(query error) %+v\n", err)
		return true, nil
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("TasksModule.SearchCreateTimeList Error: (rows close) %+v\n", err)
		}
	}()

	returnTask := make([]models.Task, 0)

	for rows.Next() {
		task := models.Task{}
		if err := rows.Scan(&task.TaskId, &task.ProjectId, &task.Name, &task.Creator, &task.Assign,
			&task.Status, &task.Deadline, &task.Description, &task.CreateDate, &task.Adopted, &task.CreatorName,
			&task.CreatorIconPath, &task.AssignName, &task.AssignIconPath); err != nil {
			log.Printf("TasksModule.SearchCreateTimeList Error: %+v\n", err)
			return true, nil
		}

		deadline, err := time.Parse("2006-01-02T15:04:05-07:00", task.Deadline)
		if err != nil {
			log.Printf("TasksModule.GetList Error: %+v\n", err)
		} else {
			deadlineFormatted := deadline.Format("2006-01-02")
			task.Deadline = deadlineFormatted
		}
		returnTask = append(returnTask, task)
	}

	return false, returnTask
}

func (self *TasksModule) Update(createDate int64, editor int, task *models.Task) (isErr bool) {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()
	tx, err := self.db.Begin()
	if err != nil {
		log.Printf("TasksModule.Update Error: (Transaction begin error) %+v\n", err)
		return true
	}

	now := time.Now().UnixNano()
	stmt, err := tx.Prepare("INSERT INTO `todo_list_history` (name, editor, status, deadline, description, createDate, updateDate, assign) VALUE (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("TasksModule.Update Error: (query error) %+v\n", err)
		if tx.Rollback() != nil {
			log.Printf("TasksModule.Update Error: (Transaction rollback error) %+v\n", err)
		}
		return true
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Printf("TasksModule.Update Error: (stmt close error) %+v\n", err)
		}
	}()
	_, err = stmt.Exec(task.Name, editor, int(task.Status), task.Deadline, task.Description, createDate, now, task.Assign)
	if err != nil {
		log.Printf("TasksModule.Update Error: (stmt exec error) %+v\n", err)
		if tx.Rollback() != nil {
			log.Printf("TasksModule.Update Error: (Transaction rollback error) %+v\n", err)
		}
		return true
	}

	_, err = tx.Exec("UPDATE `todo_list` SET adopted = ? WHERE `createDate` = ?;", now, createDate)
	if err != nil {
		log.Printf("TasksModule.Update Error: (exec error) %+v\n", err)
		if tx.Rollback() != nil {
			log.Printf("TasksModule.Update Error: (Transaction rollback error) %+v\n", err)
		}
		return true
	}
	if tx.Commit() != nil {
		log.Printf("TasksModule.Update Error: (Transaction commit error) %+v\n", err)
	}

	return false
}

func (self *TasksModule) UpdateAll(tasks []models.Task, editor int, status *models.TaskStatus, assign int, deadline *time.Time) (isErr bool) {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()
	tx, err := self.db.Begin()
	if err != nil {
		log.Printf("TasksModule.UpdateAll Error: (Transaction begin error) %+v\n", err)
		return true
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

		stmt, err := tx.Prepare("INSERT INTO `todo_list_history` (name, editor, status, deadline, description, createDate, updateDate, assign) VALUE (?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Printf("TasksModule.UpdateAll Error: (query error) %+v\n", err)
			if tx.Rollback() != nil {
				log.Printf("TasksModule.UpdateAll Error: (Transaction rollback error) %+v\n", err)
			}
			return true
		}
		defer func() {
			err := stmt.Close()
			if err != nil {
				log.Printf("TasksModule.UpdateAll Error: (stmt close error) %+v\n", err)
			}
		}()
		_, err = stmt.Exec(task.Name, editor, task.Status, updatedDeadline, task.Description, task.CreateDate, now, task.Assign)
		if err != nil {
			log.Printf("TasksModule.UpdateAll Error: (stmt exec error) %+v\n", err)
			if tx.Rollback() != nil {
				log.Printf("TasksModule.UpdateAll Error: (Transaction rollback error) %+v\n", err)
			}
			return true
		}
		task.Adopted = now

		_, err = tx.Exec("UPDATE `todo_list` SET adopted = ? WHERE `createDate` = ?;", now, task.CreateDate)
		if err != nil {
			log.Printf("TasksModule.UpdateAll Error: (exec error) %+v\n", err)
			if tx.Rollback() != nil {
				log.Printf("TasksModule.UpdateAll Error: (Transaction rollback error) %+v\n", err)
			}
			return true
		}
	}
	if tx.Commit() != nil {
		log.Printf("TasksModule.UpdateAll Error: (Transaction commit error) %+v\n", err)
	}

	return false
}

func (self *TasksModule) Delete(createDate int64) (isErr bool) {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()
	tx, err := self.db.Begin()
	if err != nil {
		log.Printf("TasksModule.Delete Error: (Transaction begin error) %+v\n", err)
		return true
	}

	_, err = tx.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	if err != nil {
		log.Printf("TasksModule.Delete Error: (exec error - FOREIGN_KEY_CHECKS=0) %+v\n", err)
		if tx.Rollback() != nil {
			log.Printf("TasksModule.Delete Error: (Transaction rollback error) %+v\n", err)
		}
		return true
	}

	_, err = tx.Exec("DELETE FROM `todo_list` WHERE `createDate` = ?;", createDate)
	if err != nil {
		log.Printf("TasksModule.Delete Error: (exec error - todo_list delete) %+v\n", err)
		if tx.Rollback() != nil {
			log.Printf("TasksModule.Delete Error: (Transaction rollback error) %+v\n", err)
		}
		return true
	}

	_, err = tx.Exec("DELETE FROM `todo_list_history` WHERE `createDate` = ?;", createDate)
	if err != nil {
		log.Printf("TasksModule.Delete Error: (exec error - todo_list_history delete) %+v\n", err)
		if tx.Rollback() != nil {
			log.Printf("TasksModule.Delete Error: (Transaction rollback error) %+v\n", err)
		}
		return true
	}

	_, err = tx.Exec("SET FOREIGN_KEY_CHECKS = 1;")
	if err != nil {
		log.Printf("TasksModule.Delete Error: (exec error - FOREIGN_KEY_CHECKS=1) %+v\n", err)
		if tx.Rollback() != nil {
			log.Printf("TasksModule.Delete Error: (Transaction rollback error) %+v\n", err)
		}
		return true
	}

	if tx.Commit() != nil {
		log.Printf("TasksModule.Delete Error: (Transaction commit error) %+v\n", err)
	}

	return false
}

func (self *TasksModule) DeleteAll(createDates []int64) (isErr bool) {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()

	deleteWhereQuery := ""
	deleteArray := make([]interface{}, 0)
	for _, createDate := range createDates {
		if len(deleteWhereQuery) > 0 {
			deleteWhereQuery += " OR "
		}
		deleteWhereQuery += " `createDate` = ? "
		deleteArray = append(deleteArray, createDate)
	}

	tx, err := self.db.Begin()
	if err != nil {
		log.Printf("TasksModule.DeleteAll Error: (Transaction begin error) %+v\n", err)
		return true
	}

	_, err = tx.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	if err != nil {
		log.Printf("TasksModule.DeleteAll Error: (exec error - FOREIGN_KEY_CHECKS=0) %+v\n", err)
		if tx.Rollback() != nil {
			log.Printf("TasksModule.DeleteAll Error: (Transaction rollback error) %+v\n", err)
		}
		return true
	}

	_, err = tx.Exec("DELETE FROM `todo_list` WHERE "+deleteWhereQuery+";", deleteArray...)
	if err != nil {
		log.Printf("TasksModule.DeleteAll Error: (exec error - todo_list delete) %+v\n", err)
		if tx.Rollback() != nil {
			log.Printf("TasksModule.DeleteAll Error: (Transaction rollback error) %+v\n", err)
		}
		return true
	}

	_, err = tx.Exec("DELETE FROM `todo_list_history` WHERE "+deleteWhereQuery+";", deleteArray...)
	if err != nil {
		log.Printf("TasksModule.DeleteAll Error: (exec error - todo_list_history delete) %+v\n", err)
		if tx.Rollback() != nil {
			log.Printf("TasksModule.DeleteAll Error: (Transaction rollback error) %+v\n", err)
		}
		return true
	}

	_, err = tx.Exec("SET FOREIGN_KEY_CHECKS = 1;")
	if err != nil {
		log.Printf("TasksModule.DeleteAll Error: (exec error - FOREIGN_KEY_CHECKS=1) %+v\n", err)
		if tx.Rollback() != nil {
			log.Printf("TasksModule.DeleteAll Error: (Transaction rollback error) %+v\n", err)
		}
		return true
	}

	if tx.Commit() != nil {
		log.Printf("TasksModule.DeleteAll Error: (Transaction commit error) %+v\n", err)
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
  tlh.name,
  creator,
  tlh.assign,
  tlh.status,
  tlh.deadline,
  tlh.description,
  todo.createDate,
  todo.adopted,
  u1.displayName,
  u1.icon_path,
  u2.displayName,
  u2.icon_path
FROM todo_list todo
  INNER JOIN todo_list_history tlh on todo.adopted = tlh.updateDate
  INNER JOIN users u1 ON u1.uuid = todo.creator
  INNER JOIN users u2 ON u2.uuid = tlh.assign
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
			&oneTask.Description, &oneTask.CreateDate, &oneTask.Adopted, &oneTask.CreatorName, &oneTask.CreatorIconPath,
			&oneTask.AssignName, &oneTask.AssignIconPath); err != nil {
			log.Printf("TasksModule.GetTasksFromUser Error: %+v\n", err)
			return true, nil
		}
		deadline, err := time.Parse("2006-01-02T15:04:05-07:00", oneTask.Deadline)
		if err != nil {
			log.Printf("TasksModule.GetList Error: %+v\n", err)
		} else {
			deadlineFormatted := deadline.Format("2006-01-02")
			oneTask.Deadline = deadlineFormatted
		}
		list = append(list, oneTask)
	}

	return false, list
}

func (t *TasksModule) GetHistoryList(createDate int64) (history []models.TaskHistory, err error) {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()

	history = make([]models.TaskHistory, 0)
	rows, err := t.db.Query(
		`SELECT
tlh.createDate, tlh.updateDate, tlh.name, tlh.editor, tlh.status, tlh.assign, tlh.deadline, tlh.description,
u1.displayName, u1.icon_path, u2.displayName, u2.icon_path
 FROM todo_list_history tlh
   INNER JOIN users u1 ON u1.uuid = tlh.editor
   INNER JOIN users u2 ON u2.uuid = tlh.assign
WHERE createDate = ? ORDER BY tlh.updateDate DESC ;`, createDate)
	if err != nil {
		log.Printf("TasksModule.GetHistoryList Error:(query error) %+v\n", err)
		return nil, err
	}
	defer func() {
		if rows.Close() != nil {
			log.Printf("TasksModule.GetHistoryList Error:(rows close error) %+v\n", err)
		}
	}()

	for rows.Next() {
		historyItem := models.TaskHistory{}
		historyTaskItem := models.TaskHistoryItem{}

		err := rows.Scan(&historyItem.CreateDate, &historyItem.UpdateDate, &historyTaskItem.Name,
			&historyTaskItem.Editor, &historyTaskItem.Status, &historyTaskItem.Assign, &historyTaskItem.Deadline,
			&historyTaskItem.Description, &historyTaskItem.EditorName, &historyTaskItem.EditorIconPath,
			&historyTaskItem.AssignName, &historyTaskItem.AssignIconPath)
		if err != nil {
			log.Printf("TasksModule.GetHistoryList Error:(scan error) %+v\n", err)
			return nil, err
		}
		historyTaskItem.CreateDate = historyItem.CreateDate
		historyTaskItem.UpdateDate = historyItem.UpdateDate
		historyItem.Task = historyTaskItem
		history = append(history, historyItem)
	}

	return history, nil
}

func (t *TasksModule) ApplyHistory(createDate, updateDate int64) error {
	t.dbLock.Lock()
	defer t.dbLock.Unlock()

	rows, err := t.db.Query("SELECT COUNT(*) FROM todo_list_history WHERE updateDate = ? AND createDate = ?;", updateDate, createDate)
	if err != nil {
		log.Printf("TasksModule.ApplyHistory Error: (exec error) %+v\n", err)
		return err
	}

	if rows.Next() {
		var count int
		err := rows.Scan(&count)
		if err != nil {
			log.Printf("TasksModule.ApplyHistory Error: (exec error) %+v\n", err)
			return err
		}
		if count == 0 {
			return errors.New("invalid createDate and updateDate.")
		}
	} else {
		return errors.New("invalid createDate and updateDate.")
	}

	tx, err := t.db.Begin()
	if err != nil {
		log.Printf("TasksModule.ApplyHistory Error: (Transaction begin error) %+v\n", err)
		return err
	}

	_, err = tx.Exec("UPDATE `todo_list` SET adopted =? WHERE `createDate` = ?;", updateDate, createDate)
	if err != nil {
		log.Printf("TasksModule.ApplyHistory Error: (exec error) %+v\n", err)
		if tx.Rollback() != nil {
			log.Printf("TasksModule.ApplyHistory Error: (Transaction rollback error) %+v\n", err)
		}
		return err
	}

	if tx.Commit() != nil {
		log.Printf("TasksModule.ApplyHistory Error: (Transaction commit error) %+v\n", err)
	}

	return nil
}
