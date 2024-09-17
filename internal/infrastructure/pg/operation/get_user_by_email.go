package operation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/javor454/newsletter-assignment/internal/application"
)

type GetUserByEmail struct {
	pgConn *sql.DB
}

type GetUserByEmailParams struct {
	Email string
}

type UserRow struct {
	ID           string
	PasswordHash string
}

func NewGetUserByEmail(pgConn *sql.DB) *GetUserByEmail {
	return &GetUserByEmail{
		pgConn: pgConn,
	}
}

func (o *GetUserByEmail) Execute(ctx context.Context, p *GetUserByEmailParams) (*UserRow, error) {
	const query = `
		SELECT id, password_hash
		FROM users
		WHERE email = $1;
	`
	var res UserRow
	if err := o.pgConn.QueryRowContext(ctx, query, p.Email).Scan(&res.ID, &res.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, application.UserNotFoundError
		}

		return nil, fmt.Errorf("failed to get user by email: %s", err.Error())
	}

	return &res, nil
}
