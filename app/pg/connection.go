package pg

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/javor454/newsletter-assignment/app/config"
)

// TODO (nice2have): prepare wrapper for readability in code

func NewConnection(conf *config.PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable", // TODO: ssl mode?
			conf.User,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.Db,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("[PG] Error connecting: %s", err.Error())
	}

	db.SetMaxOpenConns(25) // TODO: env
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("[PG] Failed to ping: %s", err.Error())
	}

	return db, nil
}
