package operation

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/javor454/newsletter-assignment/internal/application/dto"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/row"
)

// TODO (nice2have):  merge with duplicate operations?

type GetNewslettersByUserID struct {
	pgConn *sql.DB
}

type GetNewslettersByUserIDParams struct {
	UserID     string
	PageSize   int
	PageNumber int
}

func NewGetNewslettersByUserID(pgConn *sql.DB) *GetNewslettersByUserID {
	return &GetNewslettersByUserID{
		pgConn: pgConn,
	}
}

func (o *GetNewslettersByUserID) Execute(ctx context.Context, p *GetNewslettersByUserIDParams) ([]*row.Newsletter, *dto.Pagination, error) {
	const countQuery = `
        SELECT COUNT(*) 
        FROM newsletters
        WHERE user_id = $1;
    `
	const query = `
		SELECT id, public_id, name, description, created_at 
		FROM newsletters
		WHERE user_id = $1
		ORDER BY id
		LIMIT $2 OFFSET $3;
	`

	var totalItems int
	if err := o.pgConn.QueryRowContext(ctx, countQuery, p.UserID).Scan(&totalItems); err != nil {
		return nil, nil, fmt.Errorf("failed to get total count: %s", err.Error())
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(p.PageSize)))

	offset := (p.PageNumber - 1) * p.PageSize

	rows, err := o.pgConn.QueryContext(ctx, query, p.UserID, p.PageSize, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get newsletters by user id: %s", err.Error())
	}

	newsletters := make([]*row.Newsletter, 0, p.PageSize)

	for rows.Next() {
		var r row.Newsletter
		if err := rows.Scan(&r.ID, &r.PublicID, &r.Name, &r.Description, &r.CreatedAt); err != nil {
			if err := rows.Close(); err != nil {
				return nil, nil, fmt.Errorf("failed to close rows: %s", err.Error())
			}

			return nil, nil, fmt.Errorf("failed to scan row: %s", err.Error())
		}

		newsletters = append(newsletters, &r)
	}

	if err := rows.Close(); err != nil {
		return nil, nil, fmt.Errorf("failed to close rows: %s", err.Error())
	}

	return newsletters, dto.NewPagination(p.PageNumber, p.PageSize, totalPages, totalItems), nil
}
