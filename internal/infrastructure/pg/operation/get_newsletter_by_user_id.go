package operation

import (
	"context"
	"database/sql"
	"fmt"
)

type GetNewslettersByUserID struct {
	pgConn *sql.DB
}

type GetNewslettersByUserIDParams struct {
	UserID     string
	PageSize   int
	PageNumber int
}

type NewsletterRow struct {
	ID          string
	Name        string
	Description *string
}

func NewGetNewslettersByUserID(pgConn *sql.DB) *GetNewslettersByUserID {
	return &GetNewslettersByUserID{
		pgConn: pgConn,
	}
}

func (o *GetNewslettersByUserID) Execute(ctx context.Context, p *GetNewslettersByUserIDParams) ([]*NewsletterRow, error) {
	const query = `
		SELECT id, name, description 
		FROM newsletters
		WHERE user_id = $1
		ORDER BY id
		LIMIT $2 OFFSET $3;
	`

	offset := (p.PageNumber - 1) * p.PageSize

	rows, err := o.pgConn.QueryContext(ctx, query, p.UserID, p.PageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletters by user id: %s", err.Error())
	}

	newsletters := make([]*NewsletterRow, 0, p.PageSize)

	for rows.Next() {
		var r NewsletterRow
		if err := rows.Scan(&r.ID, &r.Name, &r.Description); err != nil {
			if err := rows.Close(); err != nil {
				return nil, fmt.Errorf("failed to close rows: %s", err.Error())
			}

			return nil, fmt.Errorf("failed to scan row: %s", err.Error())
		}

		newsletters = append(newsletters, &r)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %s", err.Error())
	}

	return newsletters, nil
}
