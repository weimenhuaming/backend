package utils

import (
	"encoding/base64"
	"log"

	"golang.org/x/crypto/scrypt"
)

// Bcrypt 密码加密函数
func Bcrypt(password string) string {
	const Keylen = 10
	salt := make([]byte, 8)
	salt = []byte{12, 32, 4, 66, 66, 22, 222, 11}

	HashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, Keylen)
	if err != nil {
		log.Fatal(err)
	}
	fpw := base64.StdEncoding.EncodeToString(HashPw)
	return fpw
}
