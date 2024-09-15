package domain

import "github.com/google/uuid"

type User struct {
	id       uuid.UUID
	email    *Email
	password *Password
}

func NewUser(email *Email, password *Password) *User {
	return &User{
		id:       uuid.New(),
		email:    email,
		password: password,
	}
}

func (u *User) Id() uuid.UUID {
	return u.id
}

func (u *User) Email() *Email {
	return u.email
}

func (u *User) Password() *Password {
	return u.password
}
