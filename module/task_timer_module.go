package module

import (
	"themis.mox.si/themis/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type TasksTimerModule struct {
	db      *sql.DB
	watcher *TaskTimerWatcherModule
}

func NewTasksTimerModule(db *sql.DB, watcher *TaskTimerWatcherModule) *TasksTimerModule {
	return &TasksTimerModule{
		db:      db,
		watcher: watcher,
	}
}

func (t *TasksTimerModule) TimerToggle(createDate int64, userId int) (isStart bool, err error) {
	rows, err := t.db.Query("SELECT `id`, `assign` FROM `todo_timer` WHERE `createDate` = ? AND `endDate` = 0 ORDER BY `startDate` DESC;", createDate)

	if err != nil {
		return false, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("database row close error: %+v\n", err)
		}
	}()

	if !rows.Next() {
		if err := t.createTimer(createDate, userId); err != nil {
			return false, err
		}
		return true, nil
	}

	todoTimer := &models.TodoTimer{}
	if err := rows.Scan(&todoTimer.Id, &todoTimer.Assign); err != nil {
		return false, err
	}

	if todoTimer.Assign != userId {
		return false, errors.New("invalid operation.")
	}

	if err := t.setEndTime(todoTimer.Id); err != nil {
		return false, err
	}

	return false, nil
}

func (t *TasksTimerModule) Get(taskTimerId int) (taskTimer *models.TodoTimer, err error) {
	rows, err := t.db.Query("SELECT * FROM `todo_timer` WHERE `id` = ?;", taskTimerId)
	if err != nil {
		return nil, err
	}

	taskTimer = &models.TodoTimer{}
	if !rows.Next() {
		return taskTimer, nil
	}

	err = rows.Scan(&taskTimer.Id, &taskTimer.CreateDate, &taskTimer.Assign, &taskTimer.StartDate, &taskTimer.EndDate, &taskTimer.Note)
	if err != nil {
		return nil, err
	}

	taskTimer.StartDateUnix = taskTimer.StartDate.Unix()
	taskTimer.EndDateUnix = taskTimer.EndDate.Unix()
	taskTimer.NoteString = string(taskTimer.Note)

	return taskTimer, nil
}

func (t *TasksTimerModule) Delete(taskTimerId int) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM todo_timer WHERE id = ?`, taskTimerId)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (t *TasksTimerModule) Update(taskTimerId int, taskTimer *models.TodoTimer) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE todo_timer SET startDate = ?, endDate = ?, note = ? WHERE id = ?`,
		taskTimer.StartDate,
		taskTimer.EndDate,
		taskTimer.NoteString,
		taskTimerId)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (t *TasksTimerModule) GetTaskTimerHistory(createDate int64) (history []models.TodoTimer, err error) {
	tasks := make([]models.TodoTimer, 0)

	rows, err := t.db.Query("SELECT * FROM `todo_timer` WHERE `createDate` = ? ORDER BY `startDate` DESC;", createDate)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("database row close error: %+v\n", err)
		}
	}()

	for rows.Next() {
		task := models.TodoTimer{}
		if err := rows.Scan(&task.Id, &task.CreateDate, &task.Assign, &task.StartDate, &task.EndDate, &task.Note); err != nil {
			log.Printf("sql scan error: %+v\n", err)
			continue
		}
		task.StartDateUnix = task.StartDate.Unix()
		task.EndDateUnix = task.EndDate.Unix()
		if task.EndDateUnix < 0 {
			task.EndDateUnix = 0
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t *TasksTimerModule) GetTaskTimerStatus(createDate int64) (isStart bool, err error) {
	taskTimer := models.TodoTimer{}

	rows, err := t.db.Query("SELECT * FROM `todo_timer` WHERE `createDate` = ? ORDER BY `startDate` DESC;", createDate)

	if err != nil {
		return false, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("database row close error: %+v\n", err)
		}
	}()

	if !rows.Next() {
		return false, nil
	}
	if err := rows.Scan(&taskTimer.Id, &taskTimer.CreateDate, &taskTimer.Assign, &taskTimer.StartDate, &taskTimer.EndDate, &taskTimer.Note); err != nil {
		return false, err
	}
	taskTimer.StartDateUnix = taskTimer.StartDate.Unix()
	taskTimer.EndDateUnix = taskTimer.EndDate.Unix()
	taskTimer.NoteString = string(taskTimer.Note)
	if taskTimer.EndDateUnix < 0 {
		taskTimer.EndDateUnix = 0
	}

	return taskTimer.EndDateUnix == 0, nil
}

func (t *TasksTimerModule) SearchTaskTimer(projectIds, userIds []int, startDate, endDate *time.Time) (history []models.TodoTimer, err error) {
	history = make([]models.TodoTimer, 0)

	sqlArgs := make([]interface{}, 0)
	sqlText := ""

	// users
	{
		for index, userId := range userIds {
			sqlArgs = append(sqlArgs, userId)

			if index == 0 {
				sqlText += "`assign` IN (?"
			} else {
				sqlText += ", ?"
			}

			if index == len(userIds)-1 {
				sqlText += ") "
			}
		}
	}

	//project
	{
		for index, projectId := range projectIds {
			sqlArgs = append(sqlArgs, projectId)

			if index == 0 {
				if sqlText != "" {
					sqlText += " AND "
				}

				sqlText += "`createDate` IN (SELECT `createDate` FROM `todo_list` WHERE `project` IN (?"
			} else {
				sqlText += ", ?"
			}

			if index == len(projectIds)-1 {
				sqlText += ")) "
			}
		}
	}

	//start date, end date
	{
		if (startDate != nil || endDate != nil) && sqlText != "" {
			sqlText += " AND ("
		}

		if startDate != nil {
			startDateText := startDate.Format("2006-01-02 15:04:05")
			sqlArgs = append(sqlArgs, startDateText)
			sqlText += "`startDate` >= ? "
		}

		if endDate != nil {
			endDateText := endDate.Format("2006-01-02 15:04:05")
			sqlArgs = append(sqlArgs, endDateText)

			if startDate != nil {
				sqlText += " AND "
			}

			sqlText += "`startDate` <= ? "
		}

		if startDate != nil || endDate != nil {
			sqlText += ") "
		}
	}
	sqlText = fmt.Sprintf("SELECT * FROM `todo_timer` WHERE %s ORDER BY `startDate` DESC;", sqlText)
	rows, err := t.db.Query(sqlText, sqlArgs...)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("database row close error: %+v\n", err)
		}
	}()

	for rows.Next() {
		task := models.TodoTimer{}
		if err := rows.Scan(&task.Id, &task.CreateDate, &task.Assign, &task.StartDate, &task.EndDate, &task.Note); err != nil {
			log.Printf("sql scan error: %+v\n", err)
			continue
		}
		task.StartDateUnix = task.StartDate.Unix()
		task.EndDateUnix = task.EndDate.Unix()
		task.NoteString = string(task.Note)
		if task.EndDateUnix < 0 {
			task.EndDateUnix = 0
		}
		history = append(history, task)
	}

	return history, nil
}

/**
* -------------
* private
* -------------
 */

func (t *TasksTimerModule) createTimer(createDate int64, userId int) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	todoTimer, err := tx.Prepare("INSERT INTO `todo_timer` (`createDate`, `assign`, `startDate`, `endDate`) VALUE (?, ?, NOW(), 0)")
	if err != nil {
		return err
	}

	_, err = todoTimer.Exec(createDate, userId)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return nil
}

func (t *TasksTimerModule) setEndTime(id int) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	todoTimer, err := tx.Prepare("UPDATE `todo_timer` SET `endDate` = NOW() WHERE `id` = ?;")
	if err != nil {
		return err
	}

	_, err = todoTimer.Exec(id)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return nil
}
