package operation

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/javor454/newsletter-assignment/internal/application/dto"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/row"
)

type GetNewslettersBySubscriptionEmail struct {
	pgConn *sql.DB
}

type GetNewslettersBySubscriptionEmailParams struct {
	Email      string
	PageSize   int
	PageNumber int
}

func NewGetNewslettersBySubscriptionEmail(pgConn *sql.DB) *GetNewslettersBySubscriptionEmail {
	return &GetNewslettersBySubscriptionEmail{
		pgConn: pgConn,
	}
}

func (o *GetNewslettersBySubscriptionEmail) Execute(
	ctx context.Context,
	p *GetNewslettersBySubscriptionEmailParams,
) ([]*row.Newsletter, *dto.Pagination, error) {
	const countQuery = `
        SELECT COUNT(*) as c
        FROM subscriptions s JOIN newsletters n ON n.id = s.newsletter_id 
        WHERE s.subscriber_email = $1 AND s.disabled_at IS NULL;
    `
	const query = `
		SELECT n.id, n.public_id, n.name, n.description, n.created_at
		FROM subscriptions s JOIN newsletters n ON n.id = s.newsletter_id
		WHERE s.subscriber_email = $1 AND s.disabled_at IS NULL
		ORDER BY n.id
		LIMIT $2 OFFSET $3;
	`

	var totalItems int
	if err := o.pgConn.QueryRowContext(ctx, countQuery, p.Email).Scan(&totalItems); err != nil {
		return nil, nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// TODO: could return here in case totalItems = 0

	totalPages := int(math.Ceil(float64(totalItems) / float64(p.PageSize)))

	offset := (p.PageNumber - 1) * p.PageSize

	rows, err := o.pgConn.QueryContext(ctx, query, p.Email, p.PageSize, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get newsletters by user id: %w", err)
	}

	newsletters := make([]*row.Newsletter, 0, p.PageSize)

	for rows.Next() {
		var r row.Newsletter
		if err := rows.Scan(&r.ID, &r.PublicID, &r.Name, &r.Description, &r.CreatedAt); err != nil {
			if err := rows.Close(); err != nil {
				return nil, nil, fmt.Errorf("failed to close rows: %w", err)
			}

			return nil, nil, fmt.Errorf("failed to scan row on get newsletter by subscriber: %w", err)
		}

		newsletters = append(newsletters, &r)
	}

	if err := rows.Close(); err != nil {
		return nil, nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return newsletters, dto.NewPagination(p.PageNumber, p.PageSize, totalPages, totalItems), nil
}
