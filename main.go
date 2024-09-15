package main

import (
	"fmt"

	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/javor454/newsletter-assignment/app/http-server"
	"github.com/javor454/newsletter-assignment/app/shutdown"
	"github.com/javor454/newsletter-assignment/internal"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	appConfig, err := config.CreateAppConfig()
	if err != nil {
		panic(err)
	}

	shutdownHandler := shutdown.NewHandler()
	rootCtx := shutdownHandler.CreateRootContextWithShutdown()

	httpServer := http_server.NewHttpServer()

	internal.Register()

	ginErrChan := httpServer.RunGinServer(appConfig.HttpPort)

	select {
	case err := <-ginErrChan:
		fmt.Printf("[GIN] Server error: %s\n", err.Error())
		shutdownHandler.SignalShutdown()
	case <-rootCtx.Done():
		if err := httpServer.GracefulShutdown(); err != nil {
			fmt.Printf("[GIN] Shutdown error: %s\n", err.Error())
		}

		fmt.Println("Graceful shutdown done")
	}
}
