package operation

import (
	"context"
	"database/sql"
	"fmt"
)

type CreateNewsletter struct {
	pgConn *sql.DB
}

type CreateNewsletterParams struct {
	UserID       string
	NewsletterID string
	Name         string
	Description  *string
}

func NewCreateNewsletter(pgConn *sql.DB) *CreateNewsletter {
	return &CreateNewsletter{
		pgConn: pgConn,
	}
}

func (o *CreateNewsletter) Execute(ctx context.Context, p *CreateNewsletterParams) error {
	const query = `
		INSERT INTO newsletters (user_id, id, name, description)
		VALUES ($1, $2, $3, $4);
	`
	_, err := o.pgConn.ExecContext(ctx, query, p.UserID, p.NewsletterID, p.Name, p.Description)
	if err != nil {
		return fmt.Errorf("failed to create newsletter: %s", err.Error())
	}

	return nil
}
