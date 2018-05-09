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