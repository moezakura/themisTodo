package module

import (
	"time"
	jwt "github.com/dgrijalva/jwt-go"
	"../utils"
	"log"
	"fmt"
	"os"
	"io/ioutil"
)

type SessionModule struct {
	Secret     string
	SecretByte []byte
}

type Session struct {
	Uuid int
	jwt.StandardClaims
}

func NewSessionModule() *SessionModule {
	jwtSecret := ""
	filePath := "data/secret.key"

	_, err := os.Stat(filePath)
	if err != nil {
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal("jwt secret make error!!")
		}
		defer file.Close()

		jwtSecret = utils.RandomString(128)
		file.Write(([]byte)(jwtSecret))
	} else {
		jwtTokenRaw, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal("jwt secret read error!!")
		}
		jwtSecret = string(jwtTokenRaw)
	}

	self := &SessionModule{
		jwtSecret,
		[]byte(jwtSecret),
	}

	return self
}

func (self *SessionModule) GetToken(uuid int) string {
	userSession := &Session{
		uuid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 7).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userSession)
	tokenStr, err := token.SignedString(self.SecretByte)
	if err != nil {
		log.Printf("Error SessionModule.GetToken: %+v", err)
	}
	return tokenStr
}

func (self *SessionModule) GetUuid(tokenStr string) (isExist bool, uuid int) {
	userSession := Session{}
	token, err := jwt.ParseWithClaims(tokenStr, &userSession, func(token *jwt.Token) (interface{}, error) {
		return self.SecretByte, nil
	})

	if err != nil {
		return false, 0
	}

	if claims, ok := token.Claims.(*Session); ok && token.Valid {
		return true, claims.Uuid
	} else {
		fmt.Println(err)
	}

	return false, 0
}

func (self *SessionModule) UpdateToken(tokenStr string) (isError bool, resToken string) {
	isExist, uuid := self.GetUuid(tokenStr)
	if !isExist {
		return true, ""
	}

	return false, self.GetToken(uuid)
}
