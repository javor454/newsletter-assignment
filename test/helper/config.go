package helper

import (
	"os"
	"path/filepath"

	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/sirupsen/logrus"
)

func NewAppConfig() *config.AppConfig {
	sendGridTemplateDir := "/go/src/newsletter-assignment/template"
	if Debug() {
		cwd, err := os.Getwd()
		if err != nil {
			panic("get working directory failed")
		}

		sendGridTemplateDir = filepath.Join(cwd, "./../../../template")
	}

	return &config.AppConfig{
		HttpPort:            FuncUserHttpPort,
		LogLevel:            logrus.DebugLevel,
		JwtSecret:           "0zinGG2-iDxTjLCmd4oqw29tBlhbzNITfUO-pIdyQcc=",
		CorsAllowedOrigins:  []string{"http://localhost"},
		CorsAllowedHeaders:  []string{"authorization", "content-type"},
		Timezone:            "Europe/Prague",
		SendGridApiKey:      "api-key",
		SendGridTemplateDir: sendGridTemplateDir,
		SendMail:            false,
		Host:                "http://localhost",
	}
}

func NewPostgresConfig() *config.PostgresConfig {
	host := "postgres"
	migrationsDir := "/go/src/newsletter-assignment/migration"
	if Debug() {
		host = "localhost"
		cwd, err := os.Getwd()
		if err != nil {
			panic("get working directory failed")
		}

		migrationsDir = filepath.Join(cwd, "./../../../migration")
	}

	return &config.PostgresConfig{
		User:          "newsletter-assignment",
		Password:      "pwd",
		Db:            "newsletter-assignment",
		Host:          host,
		Port:          5432,
		MigrationsDir: migrationsDir,
	}
}

func NewFirebaseConfig() *config.FirebaseConfig {
	serviceAccountFilePath := "/go/src/newsletter-assignment/service-account.json"
	if Debug() {
		cwd, err := os.Getwd()
		if err != nil {
			panic("get working directory failed")
		}
		serviceAccountFilePath = filepath.Join(cwd, "./../../../service-account.json")
	}

	return &config.FirebaseConfig{
		UseEmulator:            true,
		EmulatorHost:           "firebase-emulator:9000/?ns=strv-go-newsletter-javor-jiri",
		DatabaseURL:            "https://strv-go-newsletter-javor-jiri.firebaseio.io",
		ServiceAccountFilePath: serviceAccountFilePath,
	}
}
