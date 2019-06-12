package module

import (
	"../models"
	"database/sql"
	"time"
)

type TaskTimerWatcherModule struct {
	db *sql.DB
}

func NewTaskTimerWatcherModule(db *sql.DB) *TaskTimerWatcherModule {
	return &TaskTimerWatcherModule{
		db: db,
	}
}

func (t *TaskTimerWatcherModule) Start() {
	go func() {
		tm := time.NewTicker(10 * time.Minute)
		for {
			select {
			case <-tm.C:
				t.job()
			}
		}
		//noinspection GoUnreachableCode
		tm.Stop()
	}()
}

func (t *TaskTimerWatcherModule) job() {
	rows, err := t.db.Query("SELECT * FROM `todo_timer` WHERE `startDate` >= DATE_SUB(CURRENT_DATE(),interval 1 day) AND `endDate` = 0;")
	if err != nil {
		return
	}

	defer func() {
		if err := rows.Close(); err != nil {
			return
		}
	}()

	for rows.Next() {
		taskTimer := models.TodoTimer{}
		if err := rows.Scan(&taskTimer.Id, &taskTimer.CreateDate, &taskTimer.Assign, &taskTimer.StartDate, &taskTimer.EndDate, &taskTimer.Note); err != nil {
			continue
		}

		now := time.Now()
		if now.Day() != taskTimer.StartDate.Day() {
			tx, err := t.db.Begin()
			if err != nil {
				continue
			}

			// UPDATE
			todoTimer, err := tx.Prepare("UPDATE `todo_timer` SET `endDate` = ? WHERE `id` = ?;")
			if err != nil {
				continue
			}

			_, err = todoTimer.Exec(taskTimer.StartDate.Format("2006-01-02") + " 23:59:59", taskTimer.Id)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					continue
				}
				continue
			}

			// INSERT
			todoTimer, err = tx.Prepare("INSERT INTO `todo_timer` (`createDate`, `assign`, `startDate`, `endDate`) VALUE (?, ?, ?, 0)")
			if err != nil {
				continue
			}

			//2019-06-12 00:00:00
			tomorrow := taskTimer.StartDate.AddDate(0, 0, 1)
			_, err = todoTimer.Exec(taskTimer.CreateDate, taskTimer.Assign, tomorrow.Format("2006-01-02") + " 00:00:00")
			if err != nil {
				if err := tx.Rollback(); err != nil {
					continue
				}
				continue
			}
			if err := tx.Commit(); err != nil {
				if err := tx.Rollback(); err != nil {
					continue
				}
				continue
			}
		}
	}
}
