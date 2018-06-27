package module

import (
	"time"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"fmt"
)

type SessionModule struct {
}

type Session struct {
	Uuid int
	jwt.StandardClaims
}

func NewSessionModule() *SessionModule {
	self := &SessionModule{}
	return self
}

func (self *SessionModule) GetToken (uuid int) string {
	userSession := &Session{
		uuid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 7).Unix(),
			IssuedAt: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userSession)
	tokenStr, err := token.SignedString([]byte("foobar"))
	if err != nil {
		log.Printf("Error SessionModule.GetToken: %+v", err)
	}
	return tokenStr
}

func (self *SessionModule) GetUuid(tokenStr string) (isExist bool, uuid int) {
	userSession := Session{}
	token, err := jwt.ParseWithClaims(tokenStr, &userSession, func(token *jwt.Token) (interface{}, error) {
		return []byte("foobar"), nil
	})

	if err != nil{
		return false, 0
	}

	if claims, ok := token.Claims.(*Session); ok && token.Valid {
		return true ,claims.Uuid
	} else {
		fmt.Println(err)
	}

	return false, 0
}