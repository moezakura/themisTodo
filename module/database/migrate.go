package database

import (
	model "themis.mox.si/themis/models"
	"database/sql"
	"github.com/golang-migrate/migrate"
	"log"
)

type dbMigrate struct {
	db *sql.DB
}

func Migrate(db *sql.DB, migrate *migrate.Migrate) {
	log.Println("migration start.")
	m := &dbMigrate{
		db: db,
	}

	v, _, err := migrate.Version()
	if err != nil {
		panic(err.Error())
	}
	log.Printf("current database version: %d\n", v)

	// migrate v4
	if v < 4 {
		log.Printf("migrate %d => 4\n", v)
		err := migrate.Steps(4 - int(v))
		if err != nil {
			log.Fatalf("migrate faild (%d => 4). %+v\n", v, err)
		}
		m.ToV4()
		v = 4
	}

	if v < 6 {
		log.Printf("migrate %d => 6\n", v)
		err := migrate.Steps(6 - int(v))
		if err != nil {
			log.Fatalf("migrate faild (%d => 6). %+v\n", v, err)
		}
		m.ToV6()
	}
}

func (m *dbMigrate) ToV6() {
	rows, err := m.db.Query("SELECT `createDate` FROM `todo_list`;")
	if err != nil {
		log.Fatalf("migrate faild. %+v\n", err)
		return
	}

	createDates := make([]string, 0)
	for rows.Next() {
		var createDate string
		err := rows.Scan(&createDate)
		if err != nil {
			log.Printf("migrate skip. %+v\n", err)
			continue
		}
		createDates = append(createDates, createDate)
	}

	tx, err := m.db.Begin()
	if err != nil {
		log.Fatalf("migrate faild(transaction begin faild). %+v\n", err)
		return
	}

	for _, createDate := range createDates {
		stmt, err := tx.Prepare("UPDATE `todo_list` SET `adopted` = ? WHERE createDate = ?;")
		if err != nil {
			log.Printf("migrate skip(statement faild). %+v\n", err)
			continue
		}
		defer func() {
			err := stmt.Close()
			if err != nil {
				log.Printf("Statement close faild. %+v\n", err)
			}
		}()

		_, err = stmt.Exec(createDate, createDate)
		if err != nil {
			log.Printf("migrate skip(Statement exec faild). %+v\n", err)
			continue
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("migrate faild(transaction commit faild). %+v\n", err)
	}
}

func (m *dbMigrate) ToV4() {
	rows, err := m.db.Query("SELECT `name`, `creator`, `status`, `deadline`, `description`, `createDate`, `assign` FROM `todo_list`;")
	if err != nil {
		log.Fatalf("migrate faild. %+v\n", err)
		return
	}

	tasks := make([]model.Task, 0)
	for rows.Next() {
		task := model.Task{}
		err := rows.Scan(&task.Name, &task.Creator, &task.Status, &task.Deadline, &task.Description, &task.CreateDate, &task.Assign)
		if err != nil {
			log.Printf("migrate skip. %+v\n", err)
			continue
		}
		tasks = append(tasks, task)
	}

	tx, err := m.db.Begin()
	if err != nil {
		log.Fatalf("migrate faild(transaction begin faild). %+v\n", err)
		return
	}

	for _, task := range tasks {
		stmt, err := tx.Prepare("INSERT INTO `todo_list_history` (`name`, `editor`, `status`, `deadline`, `description`, `createDate`, `updateDate`, `assign`) VALUE (?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Printf("migrate skip.(statement faild) %+v", err)
			continue
		}

		defer func() {
			err := stmt.Close()
			if err != nil {
				log.Printf("Statement close faild. %+v\n", err)
			}
		}()

		_, err = stmt.Exec(task.Name, task.Creator, task.Status, task.Deadline, task.Description, task.CreateDate, task.CreateDate, task.Assign)
		if err != nil {
			log.Printf("migrate skip.(statement exec faild) %+v", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("migrate faild(transaction commit faild). %+v\n", err)
	}
}
