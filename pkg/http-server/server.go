package http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	engine          *gin.Engine
	srv             *http.Server
	shutdownTimeout time.Duration
}

func NewHttpServer() *HttpServer {
	// gin.SetMode(gin.ReleaseMode) PROD MODE?
	ge := gin.New()
	// TODO: cors
	// TODO: panic

	return &HttpServer{
		engine: ge,
	}
}

func (s *HttpServer) GracefulShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if s.srv == nil {
		return fmt.Errorf("[GIN] Http server not started yet")
	}

	if err := s.srv.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *HttpServer) RunGinServer() chan error {
	errChan := make(chan error, 1)

	// TODO: register in module
	s.engine.GET("/", func(c *gin.Context) {
		fmt.Println("Ping")
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	for _, v := range s.engine.Routes() {
		fmt.Printf("[GIN] Route: %s %s initialized.\n", v.Method, v.Path)
	}

	s.srv = &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Handler:      s.engine.Handler(),
	}

	go func() {
		// don't block startup with server init
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("[GIN] Server error: %w", err)
		}
	}()

	return errChan
}
