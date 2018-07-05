package utils

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type PasswordVersion int

const (
	VersionSha512 PasswordVersion = iota
	VersionHmacSha512
)

type Passworder interface {
	Hash(rawPassword string) (hashedPassword string)
	Equal(hash1, hash2 string) bool
}

func getPasswordSecret() (secret []byte, err error) {
	filePath := "data/password_secret.key"
	secretString := ""
	if _, err := os.Stat(filePath); err != nil {
		file, err := os.Create(filePath)
		if err != nil {
			return secret, errors.New("password secret make error!!")
		}
		defer file.Close()

		secretString = RandomString(128)
		file.Write(([]byte)(secretString))
	} else {
		jwtTokenRaw, err := ioutil.ReadFile(filePath)
		if err != nil {
			return secret, errors.New("password secret read error!!")
		}
		secretString = string(jwtTokenRaw)
	}
	return []byte(secretString), nil
}

func NewPassworder(version PasswordVersion) (passworder Passworder, err error) {
	switch version {
	case VersionSha512:
		return &PasswordSha512{}, nil
	case VersionHmacSha512:
		secret, err := getPasswordSecret()
		if err != nil {
			log.Fatal(err)
		}
		passwordHmac := PasswordHmacSha512{secret}
		return &passwordHmac, nil
	default:
		return nil, errors.New("Unknown password version.")
	}
}

type PasswordSha512 struct{}

func (s *PasswordSha512) Hash(rawPassword string) (hashedPassword string) {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(rawPassword)))
}

func (s *PasswordSha512) Equal(hash1, hash2 string) bool {
	return hash1 == hash2
}

type PasswordHmacSha512 struct {
	secret []byte
}

func (h *PasswordHmacSha512) Hash(rawPassword string) (hashedPassword string) {
	mac := hmac.New(sha512.New, h.secret)
	mac.Write([]byte(rawPassword))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func (h *PasswordHmacSha512) Equal(hash1, hash2 string) bool {
	byteHash1, err := hex.DecodeString(hash1)
	if err != nil {
		return false
	}
	byteHash2, err := hex.DecodeString(hash2)
	if err != nil {
		return false
	}
	return hmac.Equal(byteHash1, byteHash2)
}
