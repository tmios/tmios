package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const salt = "YOu_KONW*&_THIS!"

func GenPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("assword should not be empty")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	PasswordHash := string(passwordHash)
	return PasswordHash, nil
}

func EncryptPassword(password string) string {
	hash := sha256.Sum256(
		[]byte(fmt.Sprintf("%s-%s", password, salt)),
	)
	return fmt.Sprintf("%x", hash)
}
