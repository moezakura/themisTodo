package module

import (
	"database/sql"
	"log"
	"fmt"
)

type LoginModule struct {
	db *sql.DB
}

func NewLoginModule(db *sql.DB) *LoginModule {
	return &LoginModule{db}
}

func (self *LoginModule) IsLogin(name, password string) (error bool, authToken string) {
	var uuid int
	if err := self.db.QueryRow("SELECT count(`uuid`) FROM `users` WHERE `name` = ? AND `password` = ?;", name, password).Scan(&uuid); err != nil {
		log.Fatal(err)
		return true, ""
	}

	if uuid < 1 {
		return true, ""
	} else {
		return false, fmt.Sprintf("%d", uuid)
	}
}
