package database

import (
	model "../../models"
	"database/sql"
	"github.com/golang-migrate/migrate"
	"log"
)

type dbMigrate struct {
	db *sql.DB
}

func Migrate(db *sql.DB, migrate *migrate.Migrate) {
	m := &dbMigrate{
		db: db,
	}

	v, _, err := migrate.Version()
	if err != nil {
		panic(err.Error())
	}

	// migrate v4
	if (v < 4) {
		err := migrate.Steps(4 - int(v))
		if err != nil {
			log.Fatalf("migrate faild (%d => 4). %+v\n", v, err)
		}
		m.ToV4()
	}
}

func (m *dbMigrate) ToV4() {
	rows, err := m.db.Query("SELECT `name`, `creator`, `status`, `deadline`, `description`, `createDate`, `assign` FROM `todo_list`;")
	if err != nil {

	}

	tasks := make([]model.Task, 0)
	for rows.Next() {
		task := model.Task{}
		rows.Scan(&task.Name, &task.Creator, &task.Status, &task.Deadline, &task.Description, &task.CreateDate, &task.Assign)
		tasks = append(tasks, task)
	}

	tx, err := m.db.Begin()

	for _, task := range tasks {
		stmt, err := tx.Prepare("INSERT INTO `todo_list_history` (`name`, `editor`, `status`, `deadline`, `description`, `createDate`, `updateDate`, `assign`) VALUE (?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Printf("%+v", err)
			tx.Rollback()
		}

		defer stmt.Close()

		_, err = stmt.Exec(task.Name, task.Creator, task.Status, task.Deadline, task.Description, task.CreateDate, task.CreateDate, task.Assign)
		if err != nil {
			log.Printf("%+v", err)
			tx.Rollback()
		}
	}
	tx.Commit()
}
