package http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/javor454/newsletter-assignment/app/logger"
)

type HttpServer struct {
	engine          *gin.Engine
	lg              logger.Logger
	srv             *http.Server
	shutdownTimeout time.Duration
}

func NewHttpServer(lg logger.Logger) *HttpServer {
	// gin.SetMode(gin.ReleaseMode) PROD MODE?
	ge := gin.New()
	// TODO: cors
	// TODO: panic

	return &HttpServer{
		engine: ge,
		lg:     lg,
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

func (s *HttpServer) RunGinServer(port int) chan error {
	errChan := make(chan error, 1)

	for _, v := range s.engine.Routes() {
		s.lg.Debugf("[GIN] Route: %s %s initialized.\n", v.Method, v.Path)
	}

	s.srv = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Handler:      s.engine.Handler(),
	}

	go func() {
		// don't block startup with server init
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("[GIN] Server error: %s", err.Error())
		}
	}()

	return errChan
}

func (s *HttpServer) GetEngine() *gin.Engine {
	return s.engine
}
