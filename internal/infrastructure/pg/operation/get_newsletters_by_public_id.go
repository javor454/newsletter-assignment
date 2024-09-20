package operation

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/row"
)

// TODO (nice2have):  merge with duplicate operations or user ORM

type GetNewslettersByPublicID struct {
	pgConn *sql.DB
}

type GetNewslettersByPublicIDParams struct {
	PublicID string
}

func NewGetNewslettersByPublicID(pgConn *sql.DB) *GetNewslettersByPublicID {
	return &GetNewslettersByPublicID{
		pgConn: pgConn,
	}
}

func (o *GetNewslettersByPublicID) Execute(ctx context.Context, p *GetNewslettersByPublicIDParams) (*row.Newsletter, error) {
	const query = `
		SELECT id, public_id, name, description, created_at 
		FROM newsletters
		WHERE public_id = $1;
	`

	var newsRow row.Newsletter
	if err := o.pgConn.QueryRowContext(ctx, query, p.PublicID).Scan(
		&newsRow.ID,
		&newsRow.PublicID,
		&newsRow.Name,
		&newsRow.Description,
		&newsRow.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get newsletters by public id: %w", err)
	}

	return &newsRow, nil
}
