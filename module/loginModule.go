package module

import (
	"database/sql"
	"log"
)

type LoginModule struct {
	db *sql.DB
}

func NewLoginModule(db *sql.DB) *LoginModule {
	return &LoginModule{db}
}

func (self *LoginModule) IsLogin(name, password string) (error bool, uuid int) {
	if err := self.db.QueryRow("SELECT count(`uuid`) FROM `users` WHERE `name` = ? AND `password` = ?;", name, password).Scan(&uuid); err != nil {
		log.Fatal(err)
		return true, 0
	}

	if uuid < 1 {
		return true, 0
	} else {
		return false, uuid
	}
}
