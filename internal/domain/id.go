package domain

import (
	"github.com/google/uuid"
	"github.com/javor454/newsletter-assignment/internal/application"
)

type ID struct {
	value string
}

func NewID() *ID {
	return &ID{value: uuid.New().String()}
}

func CreateIDFromExisting(u string) (*ID, error) {
	if _, err := uuid.Parse(u); err != nil {
		return nil, application.InvalidUUIDError
	}

	return &ID{value: u}, nil
}

func (i *ID) String() string {
	return i.value
}
