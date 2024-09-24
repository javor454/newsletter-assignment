package helper

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Encrypt(subject string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(subject), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

func IsEqual(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}

	return true
}
