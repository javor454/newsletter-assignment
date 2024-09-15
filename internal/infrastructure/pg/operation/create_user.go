package operation

import "database/sql"

type CreateUser struct {
	pgConn *sql.DB
}

type CreateUserParams struct {
	Email        string
	PasswordHash stringa
}

func NewCreateUser(pgConn *sql.DB) *CreateUser {
	return &CreateUser{
		pgConn: pgConn,
	}
}

func (o *CreateUser) Execute(p *CreateUserParams) {
	const query = `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2);
	`
	o.pgConn.QueryRow(query)
}
