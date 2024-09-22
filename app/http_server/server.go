package http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	_ "github.com/javor454/newsletter-assignment/docs"
	"github.com/swaggo/swag"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/javor454/newsletter-assignment/internal/ui/http/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine          *gin.Engine
	lg              logger.Logger
	srv             *http.Server
	shutdownTimeout time.Duration
}

func NewServer(lg logger.Logger, cfg *config.AppConfig) *Server {
	lg.Debug("[GIN] Creating server...")

	gin.SetMode(gin.ReleaseMode)
	ge := gin.New()

	cf := cors.DefaultConfig()
	cf.AllowOrigins = cfg.CorsAllowedOrigins
	cf.AllowMethods = []string{"GET", "POST"}
	cf.AllowHeaders = cfg.CorsAllowedHeaders
	cf.AllowCredentials = true
	cf.MaxAge = 12 * time.Hour

	ge.Use(cors.New(cf))
	ge.Use(gin.Recovery())
	ge.Use(middleware.LoggingMiddleware(lg, []string{"/api/docs"}))

	gsc := ginSwagger.Config{
		URL:                      "doc.json",
		DocExpansion:             "list",
		InstanceName:             swag.Name,
		Title:                    "Newsletter API",
		DefaultModelsExpandDepth: 2,
		DeepLinking:              true,
		PersistAuthorization:     false,
		Oauth2DefaultClientID:    "",
	}
	ge.GET("/api/docs/*any", ginSwagger.CustomWrapHandler(&gsc, swaggerFiles.Handler))
	s := &Server{
		engine: ge,
		lg:     lg,
	}

	lg.Info("[GIN] Server created")

	return s
}

func (s *Server) GracefulShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if s.srv == nil {
		return fmt.Errorf("[GIN] Server not started yet")
	}

	if err := s.srv.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) RunGinServer(port int) chan error {
	s.lg.Debug("[GIN] Running server...")

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
			errChan <- fmt.Errorf("[GIN] Server error: %w", err)
		}
	}()
	s.lg.Info("[GIN] Server running...")

	return errChan
}

func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}
