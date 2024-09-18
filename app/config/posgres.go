package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	envPgUser     = "CONFIG_POSTGRES_USER"
	envPgPassword = "CONFIG_POSTGRES_PASSWORD"
	envPgDb       = "CONFIG_POSTGRES_DB"
	envPgHost     = "CONFIG_POSTGRES_HOST"
	envPgPort     = "CONFIG_POSTGRES_PORT"
)

type PostgresConfig struct {
	User     string
	Password string
	Db       string
	Host     string
	Port     int
}

func CreatePostgresConfig() (*PostgresConfig, error) {
	user := viper.GetString(envPgUser)
	if user == "" {
		return nil, getMissingError(envPgUser)
	}
	password := viper.GetString(envPgPassword)
	if password == "" {
		return nil, getMissingError(envPgPassword)
	}
	db := viper.GetString(envPgDb)
	if db == "" {
		return nil, getMissingError(envPgDb)
	}
	host := viper.GetString(envPgHost)
	if host == "" {
		return nil, getMissingError(envPgHost)
	}
	port := viper.GetInt(envPgPort)
	if port == 0 {
		return nil, getMissingError(envPgPort)
	}

	return &PostgresConfig{
		User:     user,
		Password: password,
		Db:       db,
		Host:     host,
		Port:     port,
	}, nil
}

func getMissingError(field string) error {
	return fmt.Errorf("missing required environment variable: %s", field)
}
