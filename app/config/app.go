package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	envHttpPort  = "CONFIG_HTTP_PORT"
	envLogLevel  = "CONFIG_LOG_LEVEL"
	envJwtSecret = "CONFIG_JWT_SECRET"
)

type AppConfig struct {
	HttpPort  int
	LogLevel  logrus.Level
	JwtSecret string
}

func CreateAppConfig() (*AppConfig, error) {
	httpPort := viper.GetInt(envHttpPort)
	if httpPort == 0 {
		return nil, getMissingError(envHttpPort)
	}
	logLevel := viper.GetString(envLogLevel)
	if logLevel == "" {
		return nil, getMissingError(envLogLevel)
	}
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}
	jwtSecret := viper.GetString(envJwtSecret)
	if jwtSecret == "" {
		return nil, getMissingError(envJwtSecret)
	}

	return &AppConfig{
		HttpPort:  httpPort,
		LogLevel:  lvl,
		JwtSecret: jwtSecret,
	}, nil
}
