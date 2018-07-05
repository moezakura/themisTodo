package module

import (
	"database/sql"
	"log"
	"github.com/gin-gonic/gin"
	"time"
	"../utils"
)

type LoginModule struct {
	db *sql.DB
}

func NewLoginModule(db *sql.DB) *LoginModule {
	return &LoginModule{db}
}

func (self *LoginModule) IsLogin(name, password string) (error bool, uuid int) {
	rows, err := self.db.Query("SELECT `uuid`, `password`, `password_version` FROM `users` WHERE `name` = ?;", name)

	if err != nil {
		log.Printf("LoginModule.IsLogin Error: %+v", err)
		return true, 0
	}

	defer rows.Close()
	rows.Next()

	var expectedPassword string
	var passwordVersion utils.PasswordVersion
	if err := rows.Scan(&uuid, &expectedPassword, &passwordVersion); err != nil {
		log.Printf("LoginModule.IsLogin Error: %+v", err)
		return true, 0
	}

	if uuid < 1 {
		return true, 0
	}

	passworder, err := utils.NewPassworder(passwordVersion)
	if err != nil {
		log.Fatal("LoginModule.IsLogin Error: %+v", err)
	}
	requestedPassword := passworder.Hash(password)
	if passworder.Equal(expectedPassword, requestedPassword) {
		return false, uuid
	} else {
		return true, 0
	}
}

func (self *LoginModule) IsLoginFromUuid(uuid int, password string) (error bool, _uuid int) {
	rows, err := self.db.Query("SELECT `uuid`, `password`, `password_version` FROM `users` WHERE `uuid` = ?;", uuid)

	if err != nil {
		log.Printf("LoginModule.IsLoginFromUuid Error: %+v", err)
		return true, 0
	}

	defer rows.Close()
	rows.Next()

	var expectedPassword string
	var passwordVersion utils.PasswordVersion
	if err := rows.Scan(&uuid, &expectedPassword, &passwordVersion); err != nil {
		log.Printf("LoginModule.IsLoginFromUuid Error: %+v", err)
		return true, 0
	}

	if uuid < 1 {
		return true, 0
	}

	passworder, err := utils.NewPassworder(passwordVersion)
	if err != nil {
		log.Fatal("LoginModule.IsLogin Error: %+v", err)
	}
	requestedPassword := passworder.Hash(password)
	if passworder.Equal(expectedPassword, requestedPassword) {
		return false, uuid
	} else {
		return true, 0
	}
}

func (self *LoginModule) GetUserId(c *gin.Context, session *SessionModule) (error bool, uuid int) {
	token, err := c.Cookie("token")

	if err != nil {
		return true, -1
	}

	exist, userUuid := session.GetUuid(token)
	if !exist {
		return true, -1
	}

	tokenUpdateTime := time.Now().AddDate(0, 0, 2)
	if session.GetExpires(token) < tokenUpdateTime.Unix() {
		limitSec := 30 * 24 * 60 * 60
		isErr, authToken := session.UpdateToken(token)

		if isErr {
			return false, userUuid
		}

		c.SetCookie("token", authToken, limitSec, "/", "", false, true)
	}

	return false, userUuid
}
