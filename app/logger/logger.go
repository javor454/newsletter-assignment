package logger

import (
	"os"

	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	logrus.FieldLogger
}

func NewLogger(appConfig *config.AppConfig) Logger {
	logger := logrus.New()
	logger.Level = appConfig.LogLevel
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	logger.Debugln("[LOGGER] Creating...")

	return logger
}
