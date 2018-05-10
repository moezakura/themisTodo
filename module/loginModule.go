package module

import (
	"database/sql"
	"log"
	"github.com/gin-gonic/gin"
)

type LoginModule struct {
	db *sql.DB
}

func NewLoginModule(db *sql.DB) *LoginModule {
	return &LoginModule{db}
}

func (self *LoginModule) IsLogin(name, password string) (error bool, uuid int) {
	if err := self.db.QueryRow("SELECT `uuid` FROM `users` WHERE `name` = ? AND `password` = ?;", name, password).Scan(&uuid); err != nil {
		log.Printf("LoginModule.IsLogin Error: %+v", err)
		return true, 0
	}

	if uuid < 1 {
		return true, 0
	} else {
		return false, uuid
	}
}

func (self *LoginModule) GetUserId(c *gin.Context, session *SessionModule) (error bool, uuid int){
	token, err := c.Cookie("token")

	if err != nil {
		return true, -1
	}

	exist, userUuid := session.GetUuid(token)
	if !exist {
		return true, -1
	}

	return false, userUuid
}