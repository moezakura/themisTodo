package module

import (
	"../models"
	"database/sql"
	"errors"
	"github.com/labstack/gommon/log"
)

type TasksTimerModule struct {
	db *sql.DB
}

func NewTasksTimerModule(db *sql.DB) *TasksTimerModule {
	return &TasksTimerModule{db: db}
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
