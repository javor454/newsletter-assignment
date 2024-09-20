package operation

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type UpdateUnsentEmailJobs struct {
	pgConn *sql.DB
}

type UpdateUnsentEmailJobsParams struct {
	JobIDs []string
}

func NewUpdateUnsentEmailJobs(pgConn *sql.DB) *UpdateUnsentEmailJobs {
	return &UpdateUnsentEmailJobs{
		pgConn: pgConn,
	}
}

func (u *UpdateUnsentEmailJobs) Execute(ctx context.Context, p *UpdateUnsentEmailJobsParams) error {
	const query = "UPDATE email_jobs SET sent = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id = ANY($1);"

	_, err := u.pgConn.ExecContext(ctx, query, pq.Array(p.JobIDs))
	if err != nil {
		return fmt.Errorf("failed to execute update unsent email jobs: %w", err)
	}

	return nil
}
