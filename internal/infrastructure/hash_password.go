package infrastructure

import "golang.org/x/crypto/bcrypt"

type HashPassword struct{}

func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash), err
}
