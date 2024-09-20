package operation

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/javor454/newsletter-assignment/internal/application"
)

type CreateNewsletter struct {
	pgConn *sql.DB
}

type CreateNewsletterParams struct {
	UserID      string
	ID          string
	PublicID    string
	Name        string
	Description *string
	CreatedAt   time.Time
}

func NewCreateNewsletter(pgConn *sql.DB) *CreateNewsletter {
	return &CreateNewsletter{
		pgConn: pgConn,
	}
}

func (o *CreateNewsletter) Execute(ctx context.Context, p *CreateNewsletterParams) error {
	const (
		unknownUserConstraint = "newsletters_user_id_fkey"
		query                 = `
			INSERT INTO newsletters (user_id, id, public_id, name, description, created_at)
			VALUES ($1, $2, $3, $4, $5, $6);
		`
	)
	_, err := o.pgConn.ExecContext(ctx, query, p.UserID, p.ID, p.PublicID, p.Name, p.Description, p.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), unknownUserConstraint) {
			return application.UnknownUserError
		}

		return fmt.Errorf("failed to create newsletter: %w", err)
	}

	return nil
}
