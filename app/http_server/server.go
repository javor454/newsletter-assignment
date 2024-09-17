package http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/javor454/newsletter-assignment/app/logger"
)

type Server struct {
	engine          *gin.Engine
	lg              logger.Logger
	srv             *http.Server
	shutdownTimeout time.Duration
}

func NewServer(lg logger.Logger, cfg *config.AppConfig) *Server {
	gin.SetMode(gin.DebugMode) // TODO: PROD MODE?
	ge := gin.New()

	cf := cors.DefaultConfig()
	cf.AllowOrigins = cfg.CorsAllowedOrigins
	cf.AllowMethods = []string{"GET", "POST"}
	cf.AllowHeaders = cfg.CorsAllowedHeaders
	cf.AllowCredentials = true
	cf.MaxAge = 12 * time.Hour
	ge.Use(cors.New(cf))

	// TODO: panic middleware

	return &Server{
		engine: ge,
		lg:     lg,
	}
}

func (s *Server) GracefulShutdown() error {
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

func (s *Server) RunGinServer(port int) chan error {
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

func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}
