package operation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/javor454/newsletter-assignment/internal/application"
)

type GetNewsletterIDByPublicID struct {
	pgConn *sql.DB
}

type GetNewsletterIDByPublicIDParams struct {
	PublicID string
}

type NewsletterIDRow struct {
	ID string
}

func NewGetNewsletterIDByPublicID(pgConn *sql.DB) *GetNewsletterIDByPublicID {
	return &GetNewsletterIDByPublicID{
		pgConn: pgConn,
	}
}

func (o *GetNewsletterIDByPublicID) Execute(ctx context.Context, p *GetNewsletterIDByPublicIDParams) (*NewsletterIDRow, error) {
	const query = "SELECT id FROM newsletters WHERE public_id = $1;"

	var res NewsletterIDRow
	if err := o.pgConn.QueryRowContext(ctx, query, p.PublicID).Scan(&res.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, application.NewsletterNotFoundError
		}

		return nil, fmt.Errorf("failed to get id by public id: %s", err.Error())
	}

	return &res, nil
}
