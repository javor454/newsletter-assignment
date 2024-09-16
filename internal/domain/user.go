package domain

type User struct {
	id       *ID
	email    *Email
	password *Password
}

func NewUser(email *Email, password *Password) *User {
	return &User{
		id:       NewID(),
		email:    email,
		password: password,
	}
}

func CreateUserFromExisting(id *ID, email *Email, password *Password) *User {
	return &User{
		id:       id,
		email:    email,
		password: password,
	}
}

func (u *User) Id() *ID {
	return u.id
}

func (u *User) Email() *Email {
	return u.email
}

func (u *User) Password() *Password {
	return u.password
}
