package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	envHttpPort           = "CONFIG_HTTP_PORT"
	envLogLevel           = "CONFIG_LOG_LEVEL"
	envJwtSecret          = "CONFIG_JWT_SECRET"
	envCorsAllowedOrigins = "CONFIG_CORS_ALLOWED_ORIGINS"
	envCorsAllowedHeaders = "CONFIG_CORS_ALLOWED_HEADERS"
	envTimezone           = "CONFIG_TIMEZONE"
	envSendGridApiKey     = "CONFIG_SENDGRID_API_KEY"
	envSendMail           = "CONFIG_SEND_MAIL"
	envHost               = "CONFIG_HOST"
)

type AppConfig struct {
	HttpPort           int
	LogLevel           logrus.Level
	JwtSecret          string
	CorsAllowedOrigins []string
	CorsAllowedHeaders []string
	Timezone           string
	SendGridApiKey     string
	SendMail           bool
	Host               string
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
	corsAllowedOrigins := viper.GetStringSlice(envCorsAllowedOrigins)
	if len(corsAllowedOrigins) == 0 {
		return nil, getMissingError(envCorsAllowedOrigins)
	}
	corsAllowedHeaders := viper.GetStringSlice(envCorsAllowedHeaders)
	if len(corsAllowedHeaders) == 0 {
		return nil, getMissingError(envCorsAllowedHeaders)
	}
	timezone := viper.GetString(envTimezone)
	if timezone == "" {
		return nil, getMissingError(envTimezone)
	}
	sendGridApiKey := viper.GetString(envSendGridApiKey)
	if sendGridApiKey == "" {
		return nil, getMissingError(envSendGridApiKey)
	}
	sendMail := viper.GetBool(envSendMail)
	host := viper.GetString(envHost)
	if host == "" {
		return nil, getMissingError(envHost)
	}

	return &AppConfig{
		HttpPort:           httpPort,
		LogLevel:           lvl,
		JwtSecret:          jwtSecret,
		CorsAllowedOrigins: corsAllowedOrigins,
		CorsAllowedHeaders: corsAllowedHeaders,
		Timezone:           timezone,
		SendGridApiKey:     sendGridApiKey,
		SendMail:           sendMail,
		Host:               host,
	}, nil
}
