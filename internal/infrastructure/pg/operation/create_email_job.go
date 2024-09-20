package operation

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/row"
)

type CreateEmailJobParams struct {
	ID     string
	Type   row.MailType
	Params []byte
}

func CreateEmailJobTx(ctx context.Context, tx *sql.Tx, p *CreateEmailJobParams) error {
	const (
		query = `
			INSERT INTO email_jobs (id, message_type, params)
			VALUES ($1, $2, $3);
		`
	)

	_, err := tx.ExecContext(ctx, query, p.ID, p.Type, p.Params)
	if err != nil {
		return fmt.Errorf("failed to enqueue mail job: %w", err)
	}

	return nil
}
