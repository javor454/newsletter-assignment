package operation

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/row"
)

// TODO (nice2have): make this generic for all message types

type GetUnsentSubscribedEmailJobsOperation struct {
	pgConn *sql.DB
}

func NewGetUnsentSubscribedEmailJobsOperation(pgConn *sql.DB) *GetUnsentSubscribedEmailJobsOperation {
	return &GetUnsentSubscribedEmailJobsOperation{
		pgConn: pgConn,
	}
}

func (o *GetUnsentSubscribedEmailJobsOperation) Execute(ctx context.Context, maxJobs int) ([]*row.EmailJob, error) {
	rows, err := o.pgConn.QueryContext(
		ctx,
		fmt.Sprintf(
			"SELECT id, message_type, params FROM email_jobs WHERE sent = FALSE AND message_type = '%s' LIMIT %d;",
			row.SubscriptionType,
			maxJobs,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get unsent email jobs: %w", err)
	}

	jobs := make([]*row.EmailJob, 0, 100)

	for rows.Next() {
		var r row.EmailJob
		if err := rows.Scan(&r.ID, &r.Type, &r.Params); err != nil {
			if err := rows.Close(); err != nil {
				return nil, fmt.Errorf("failed to close rows: %w", err)
			}

			return nil, fmt.Errorf("failed to scan row on get unsent jobs: %w", err)
		}

		jobs = append(jobs, &r)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close rows: %w", err)
	}

	return jobs, nil
}
