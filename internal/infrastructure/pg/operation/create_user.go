package operation

import (
	"context"
	"database/sql"
	"fmt"
)

type CreateUser struct {
	pgConn *sql.DB
}

type CreateUserParams struct {
	ID           string
	Email        string
	PasswordHash string
}

func NewCreateUser(pgConn *sql.DB) *CreateUser {
	return &CreateUser{
		pgConn: pgConn,
	}
}

// TODO: test duplicit save
func (o *CreateUser) Execute(ctx context.Context, p *CreateUserParams) error {
	const query = `
		INSERT INTO users (id, email, password_hash)
		VALUES ($1, $2, $3);
	`
	_, err := o.pgConn.ExecContext(ctx, query, p.ID, p.Email, p.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to create user: %s", err.Error())
	}

	return nil
}
