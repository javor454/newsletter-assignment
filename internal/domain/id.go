package domain

import "github.com/google/uuid"

type ID struct {
	value string
}

func NewID() *ID {
	return &ID{value: uuid.New().String()}
}

func CreateIDFromExisting(uuid string) *ID {
	return &ID{value: uuid}
}

func (i *ID) String() string {
	return i.value
}
