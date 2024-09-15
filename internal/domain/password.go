package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hash string
}

func NewPassword(value string) (*Password, error) {
	h, err := hashBcrypt(value)
	if err != nil {
		return nil, err
	}

	return &Password{hash: h}, nil
}

func (p *Password) IsEqual(value string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(value)); err != nil {
		return false
	}

	return true
}

func (p *Password) String() string {
	return p.hash
}

func hashBcrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %s", err.Error())
	}

	return string(hash), nil
}
