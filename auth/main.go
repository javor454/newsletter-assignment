package main

import (
	"fmt"

	"github.com/javor454/newsletter-assignment/auth/internal"
	"github.com/javor454/newsletter-assignment/pkg/http-server"
	"github.com/javor454/newsletter-assignment/pkg/shutdown"
)

func main() {
	shutdownHandler := shutdown.NewHandler()
	rootCtx := shutdownHandler.CreateRootContextWithShutdown()

	httpServer := http_server.NewHttpServer()

	internal.Register()

	ginErrChan := httpServer.RunGinServer()

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
