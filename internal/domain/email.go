package domain

import (
	"fmt"
	"regexp"
)

type Email struct {
	value string
}

func NewEmail(value string) (*Email, error) {
	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !regex.MatchString(value) {
		return nil, fmt.Errorf("invalid email format")
	}

	return &Email{value: value}, nil
}

func (e *Email) String() string {
	return e.value
}
