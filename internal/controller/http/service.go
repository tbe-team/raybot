package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	httphandler "github.com/tbe-team/raybot/internal/controller/http/handler"
	"github.com/tbe-team/raybot/internal/controller/http/middleware"
	"github.com/tbe-team/raybot/internal/controller/http/oas/gen"
	"github.com/tbe-team/raybot/internal/controller/http/swagger"
	"github.com/tbe-team/raybot/internal/service"
)

const (
	DefaultPort = 3000
)

type Config struct {
	EnableSwagger bool `yaml:"enable_swagger"`
}

func (c *Config) Validate() error {
	return nil
}

type CleanupFunc func(ctx context.Context) error

//nolint:revive
type HTTPService struct {
	cfg Config

	service service.Service
	log     *slog.Logger
}

func NewHTTPService(cfg Config, service service.Service, log *slog.Logger) (*HTTPService, error) {
	return &HTTPService{
		cfg:     cfg,
		service: service,
		log:     log.With(slog.String("service", "HTTPService")),
	}, nil
}

func (s HTTPService) Run() (CleanupFunc, error) {
	r := chi.NewRouter()

	s.registerMiddlewares(r)
	if s.cfg.EnableSwagger {
		s.registerSwaggerHandler(r)
	}
	s.RegisterUIHandler(r)
	s.RegisterAPIHandlers(r)

	return s.RunWithServer(r)
}

func (s HTTPService) RunWithServer(r chi.Router) (CleanupFunc, error) {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", DefaultPort),
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		s.log.Info(fmt.Sprintf("HTTP server is listening on port %d", DefaultPort))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("HTTP server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	cleanup := func(ctx context.Context) error {
		s.log.Debug("HTTP server is shutting down")
		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown HTTP server: %w", err)
		}
		s.log.Debug("HTTP server shut down complete")
		return nil
	}

	return cleanup, nil
}

func (s HTTPService) RegisterAPIHandlers(r chi.Router) {
	apiHandler := httphandler.NewAPIHandler(s.service)
	strictAPIHandler := gen.NewStrictHandlerWithOptions(
		apiHandler,
		[]gen.StrictMiddlewareFunc{},
		gen.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  s.handleRequestError,
			ResponseErrorHandlerFunc: s.handleResponseError,
		},
	)
	gen.HandlerFromMuxWithBaseURL(strictAPIHandler, r, "/api/v1")
}

func (s HTTPService) registerMiddlewares(r chi.Router) {
	r.Use(middleware.Logging)
	r.Use(middleware.Recoverer)
}

func (s HTTPService) registerSwaggerHandler(r chi.Router) {
	swagger.Register(r, "/docs/openapi.yml")
}
