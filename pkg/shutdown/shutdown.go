package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Handler struct {
	shutdownChan chan os.Signal
}

func NewHandler() *Handler {
	return &Handler{
		shutdownChan: make(chan os.Signal, 1),
	}
}

// CreateRootContextWithShutdown Creates a context which is cancelled on SIGINT or SIGTERM.
func (s *Handler) CreateRootContextWithShutdown() context.Context {
	fmt.Println("Creating root context")
	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(s.shutdownChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-s.shutdownChan
		fmt.Println("Received shutdown signal, shutting down gracefully...")
		cancel()
	}()

	return ctx
}

func (s *Handler) SignalShutdown() {
	if s.shutdownChan != nil {
		s.shutdownChan <- syscall.SIGINT
	}
}
