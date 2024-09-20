package crypt

import (
	"fmt"

	"github.com/javor454/newsletter-assignment/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHash struct {
	value string
}

func NewBcryptHashFromPassword(password *domain.Password) (*BcryptHash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password.String()), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return &BcryptHash{
		value: string(hash),
	}, nil
}

func CreateHashFromExisting(hash string) *BcryptHash {
	return &BcryptHash{value: hash}
}

func (h *BcryptHash) IsEqual(password *domain.Password) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(h.value), []byte(password.String())); err != nil {
		return false
	}

	return true
}

func (h *BcryptHash) String() string {
	return h.value
}
