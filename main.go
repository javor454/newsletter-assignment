package main

import (
	"time"

	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/javor454/newsletter-assignment/app/firebase"
	"github.com/javor454/newsletter-assignment/app/http_server"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/javor454/newsletter-assignment/app/pg"
	"github.com/javor454/newsletter-assignment/app/sendgrid"
	"github.com/javor454/newsletter-assignment/app/shutdown"
	"github.com/javor454/newsletter-assignment/internal"
	"github.com/spf13/viper"
)

// @title			Newsletter assignment
// @version		1.0
// @description	Newsletter assignment for STRV.
// @contact.email	javornicky.jiri@gmail.com
func main() {
	shutdownHandler := shutdown.NewHandler()

	rootCtx := shutdownHandler.CreateRootContextWithShutdown()

	viper.AutomaticEnv()
	appConfig, err := config.NewAppConfig()
	if err != nil {
		panic("[CONFIG] failed to load: " + err.Error())
	}
	pgConfig, err := config.NewPostgresConfig()
	if err != nil {
		panic("[CONFIG] failed to load: " + err.Error())
	}
	fbConfig, err := config.NewFirebaseConfig()
	if err != nil {
		panic("[CONFIG] failed to load: " + err.Error())
	}

	location, err := time.LoadLocation(appConfig.Timezone)
	if err != nil {
		panic("failed to load timezone")
	}
	time.Local = location

	lg := logger.NewLogger(appConfig)

	pgConn, err := pg.NewConnection(lg, pgConfig)
	if err != nil {
		panic("[PG] failed to connect: " + err.Error())
	}

	fbClient, err := firebase.NewClient(lg, rootCtx, fbConfig)
	if err != nil {
		panic("[FIREBASE] failed to connect: " + err.Error())
	}

	mailClient := sendgrid.NewClient(lg, appConfig)
	if err != nil {
		panic("[SENDGRID] failed to create client: " + err.Error())
	}

	if err := pg.MigrationsUp(lg, pgConfig, pgConn); err != nil {
		panic("[MIGRATIONS] failed to run: " + err.Error())
	}

	httpServer := http_server.NewServer(lg, appConfig)

	internal.RegisterDependencies(rootCtx, lg, appConfig, pgConn, httpServer, mailClient, fbClient)

	ginErrChan := httpServer.RunGinServer(appConfig.HttpPort)

	select {
	case err := <-ginErrChan:
		lg.Errorf("[GIN] Server error: %s\n", err.Error())
		shutdownHandler.SignalShutdown()
	case <-rootCtx.Done():
		lg.Info("Received signal, shutting down with grace...")
		if err := httpServer.GracefulShutdown(); err != nil {
			lg.WithError(err).Fatalf("[GIN] Shutdown error.")
		}
		if err := pgConn.Close(); err != nil {
			lg.WithError(err).Fatalf("[PG] Shutdown error.")
		}

		lg.Info("Graceful shutdown done")
	}
}
