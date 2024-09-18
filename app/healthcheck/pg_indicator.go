package healthcheck

import (
	"context"
	"database/sql"
	"time"
)

type PgIndicator struct {
	conn               *sql.DB
	healthCheckTimeout time.Duration
}

func NewPgIndicator(conn *sql.DB, healthCheckTimeout time.Duration) *PgIndicator {
	return &PgIndicator{conn: conn, healthCheckTimeout: healthCheckTimeout}
}

func (p *PgIndicator) Name() string {
	return "postgres"
}

func (p *PgIndicator) Check(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, p.healthCheckTimeout)
	defer cancel()

	if err := p.conn.PingContext(ctx); err != nil {
		return false
	}

	return true
}
