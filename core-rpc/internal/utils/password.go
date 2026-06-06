package utils

import (
	"encoding/base64"
	"log"

	"golang.org/x/crypto/scrypt"
)

const keyLen = 10

var salt = []byte{12, 32, 4, 66, 66, 22, 222, 11}

func HashPassword(password string) string {
	hashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, keyLen)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(hashPw)
}

func VerifyPassword(password, hashed string) bool {
	if password == "" || hashed == "" {
		return false
	}
	hashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, keyLen)
	if err != nil {
		return false
	}
	return base64.StdEncoding.EncodeToString(hashPw) == hashed
}
