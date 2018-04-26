package utils

import (
	"fmt"
	"crypto/sha512"
)

func SHA512(plain string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(plain)))
}
