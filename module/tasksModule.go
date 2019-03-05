package module

import (
	"../models"
	"database/sql"
	"log"
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

	rows, err := self.db.Query(`SELECT id, tlh.name, creator, tlh.assign, tlh.status, tlh.deadline, tlh.description, todo.createDate, u1.displayName, u1.icon_path, u2.displayName, u2.icon_path FROM todo_list todo
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
			&listOne.Status, &listOne.Deadline, &listOne.Description,
			&listOne.CreateDate, &listOne.CreatorName, &listOne.CreatorIconPath, &listOne.AssignName,
			&listOne.AssignIconPath); err != nil {
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
  tlh.name,
  creator,
  tlh.assign,
  tlh.status,
  tlh.deadline,
  tlh.description,
  todo.createDate,
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
		&returnTask.Description, &returnTask.CreateDate, &returnTask.CreatorName, &returnTask.CreatorIconPath,
		&returnTask.AssignName, &returnTask.AssignIconPath); err != nil {
		log.Printf("TasksModule.Get Error: %+v\n", err)
		return true, nil
	}

	return false, returnTask
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

	if !rows.Next() {
		return true, nil
	}

	returnTask := make([]models.Task, 0)

	for rows.Next() {
		task := models.Task{}
		var updateDate int64
		if err := rows.Scan(&task.TaskId, &task.ProjectId, &task.Name, &task.Creator, &task.Assign,
			&task.Status, &task.Deadline, &task.Description, &task.CreateDate, &task.CreatorName,
			&task.CreatorIconPath, &task.AssignName, &task.AssignIconPath, &updateDate); err != nil {
			log.Printf("TasksModule.Search Error:(scan error) %+v\n", err)
			return true, nil
		}
		if updateDate == 0 {
			continue
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
		queryString += " createDate = ? "
		queryArray = append(queryArray, value)
	}

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
  u1.icon_path,
  u2.displayName,
  u2.icon_path
FROM todo_list todo
  INNER JOIN users u1 ON u1.uuid = todo.creator
  INNER JOIN users u2 ON u2.uuid = todo.assign
WHERE `+queryString+";", queryArray...)

	if err != nil {
		return true, nil
	}

	defer rows.Close()

	returnTask := make([]models.Task, 0)

	for rows.Next() {
		task := models.Task{}
		if err := rows.Scan(&task.TaskId, &task.ProjectId, &task.Name, &task.Creator, &task.Assign,
			&task.Status, &task.Deadline, &task.Description, &task.CreateDate, &task.CreatorName,
			&task.CreatorIconPath, &task.AssignName, &task.AssignIconPath); err != nil {
			log.Printf("TasksModule.SearchCreateTimeList Error: %+v\n", err)
			return true, nil
		}
		returnTask = append(returnTask, task)
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

func (self *TasksModule) UpdateAll(tasks []models.Task, status *models.TaskStatus, assign int, deadline *time.Time) (isErr bool) {
	self.dbLock.Lock()
	defer self.dbLock.Unlock()

	updateSetQuery := ""
	updateArray := make([]interface{}, 0)
	{
		if status != nil {
			updateSetQuery += "`status` = ? "
			updateArray = append(updateArray, status)
		}
		if assign > 0 {
			if len(updateSetQuery) > 0 {
				updateSetQuery += ","
			}
			updateSetQuery += "`assign` = ? "
			updateArray = append(updateArray, assign)
		}
		if deadline != nil {
			if len(updateSetQuery) > 0 {
				updateSetQuery += ","
			}
			updateSetQuery += "`assign` = ? "
			updateArray = append(updateArray, deadline)
		}
	}

	updateWhereQuery := ""
	for _, task := range tasks {
		if len(updateWhereQuery) > 0 {
			updateWhereQuery += " OR "
		}
		updateWhereQuery += " `createDate` = ? "
		updateArray = append(updateArray, task.CreateDate)
	}

	_, err := self.db.Exec("UPDATE `todo_list` SET "+updateSetQuery+" WHERE "+updateWhereQuery+";",
		updateArray...)

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

	_, err := self.db.Exec("DELETE FROM `todo_list` WHERE "+deleteWhereQuery+";", deleteArray...)

	if err != nil {
		log.Printf("TasksModule.DeleteAll Error: %+v\n", err)
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
  u1.icon_path,
  u2.displayName,
  u2.icon_path
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
			&oneTask.Description, &oneTask.CreateDate, &oneTask.CreatorName, &oneTask.CreatorIconPath,
			&oneTask.AssignName, &oneTask.AssignIconPath); err != nil {
			log.Printf("TasksModule.GetTasksFromUser Error: %+v\n", err)
			return true, nil
		}
		list = append(list, oneTask)
	}

	return false, list
}
