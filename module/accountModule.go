package module

import (
	"database/sql"
	"log"
	"../utils"
)

type AccountModule struct {
	db *sql.DB
}

func NewAccountModule(db *sql.DB) *AccountModule {
	return &AccountModule{db}
}

func (self *AccountModule) Add(name, password string) bool {
	stmt, err := self.db.Prepare("INSERT INTO `users` (`displayName`, `name`, `password`) VALUES(?, ?, ?);")

	if err != nil {
		log.Printf("ProjectsModule.AddUser Error: %+v", err)
		return true
	}

	defer stmt.Close()

	_, err = stmt.Exec(name, name, utils.SHA512(password))
	if err != nil {
		log.Printf("ProjectsModule.AddUser Error: %+v", err)
		return true
	}

	return false
}

func (self *AccountModule) Get(name string) int {
	rows, err := self.db.Query("SELECT `uuid` FROM `users` WHERE `name` = ?;", name)

	if err != nil {
		log.Printf("ProjectsModule.AddUser Error: %+v", err)
		return 0
	}

	defer rows.Close()

	var userId int

	if rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			log.Printf("ProjectsModule.AddUser Error: %+v", err)
			return 0
		}

		return userId
	}

	return 0
}
