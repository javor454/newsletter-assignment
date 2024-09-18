package main

import (
	"log"

	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/javor454/newsletter-assignment/app/http_server"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/javor454/newsletter-assignment/app/pg"
	"github.com/javor454/newsletter-assignment/app/shutdown"
	"github.com/javor454/newsletter-assignment/internal"
	"github.com/spf13/viper"
)

//	@title			Newsletter assignment
//	@version		1.0
//	@description	Newsletter assignment for STRV.
//	@contact.email	javornicky.jiri@gmail.com
func main() {
	viper.AutomaticEnv()
	appConfig, err := config.CreateAppConfig()
	if err != nil {
		panic(err)
	}
	pgConfig, err := config.CreatePostgresConfig()
	if err != nil {
		panic(err)
	}

	lg := logger.NewLogger(appConfig)

	lg.Debug("[PG] Connecting...")
	pgConn, err := pg.NewConnection(pgConfig)
	if err != nil {
		lg.Fatal(err)
	}
	lg.Info("[PG] Connected")

	lg.Debug("[MIGRATIONS] Starting up...")
	if err := pg.MigrationsUp(pgConn); err != nil {
		log.Fatal(err)
	}
	lg.Info("[MIGRATIONS] Done")

	shutdownHandler := shutdown.NewHandler()

	lg.Debug("Creating root context")
	rootCtx := shutdownHandler.CreateRootContextWithShutdown()

	lg.Debug("[HTTP] Creating server...")
	httpServer := http_server.NewServer(lg, appConfig)
	lg.Debug("[HTTP] Server created")

	internal.RegisterDependencies(httpServer, lg, pgConn, appConfig)

	lg.Debug("[HTTP] Running server...")
	ginErrChan := httpServer.RunGinServer(appConfig.HttpPort)
	lg.Info("[HTTP] Server running...")

	select {
	case err := <-ginErrChan:
		lg.Errorf("[HTTP] Server error: %s\n", err.Error())
		shutdownHandler.SignalShutdown()
	case <-rootCtx.Done():
		lg.Info("Received signal, shutting down with grace...")
		if err := httpServer.GracefulShutdown(); err != nil {
			lg.WithError(err).Fatalf("[HTTP] Shutdown error.")
		}
		if err := pgConn.Close(); err != nil {
			lg.WithError(err).Fatalf("[PG] Shutdown error.")
		}

		lg.Info("Graceful shutdown done")
	}
}
