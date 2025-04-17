package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	"github.com/tbe-team/raybot/internal/handlers/http/middleware"
	"github.com/tbe-team/raybot/internal/handlers/http/swagger"
	"github.com/tbe-team/raybot/internal/services/apperrorcode"
	"github.com/tbe-team/raybot/internal/services/command"
	configsvc "github.com/tbe-team/raybot/internal/services/config"
	"github.com/tbe-team/raybot/internal/services/dashboarddata"
	"github.com/tbe-team/raybot/internal/services/peripheral"
	"github.com/tbe-team/raybot/internal/services/system"
)

type Service struct {
	cfg config.HTTP
	log *slog.Logger

	configService        configsvc.Service
	systemService        system.Service
	dashboardDataService dashboarddata.Service
	peripheralService    peripheral.Service
	commandService       command.Service
	apperrorcodeService  apperrorcode.Service
}

type CleanupFunc func(ctx context.Context) error

func New(
	cfg config.HTTP,
	log *slog.Logger,
	configService configsvc.Service,
	systemService system.Service,
	dashboardDataService dashboarddata.Service,
	peripheralService peripheral.Service,
	commandService command.Service,
	apperrorcodeService apperrorcode.Service,
) *Service {
	return &Service{
		cfg:                  cfg,
		log:                  log.With("service", "http"),
		configService:        configService,
		systemService:        systemService,
		dashboardDataService: dashboardDataService,
		peripheralService:    peripheralService,
		commandService:       commandService,
		apperrorcodeService:  apperrorcodeService,
	}
}

func (s *Service) Run() (CleanupFunc, error) {
	r := chi.NewRouter()
	s.RegisterMiddlewares(r)

	if s.cfg.Swagger {
		swagger.Register(r)
	}

	s.RegisterUIHandler(r)
	s.RegisterHandlers(r)

	return s.RunWithServer(r)
}

func (s *Service) RunWithServer(handler http.Handler) (CleanupFunc, error) {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.cfg.Port),
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("failed to start HTTP server", slog.Any("error", err))
		}
	}()

	return func(ctx context.Context) error {
		return srv.Shutdown(ctx)
	}, nil
}

func (s *Service) RegisterMiddlewares(r chi.Router) {
	r.Use(middleware.Recoverer(s.log))
	r.Use(middleware.Logging(s.log))
	r.Use(middleware.Cors())
}

func (s *Service) RegisterHandlers(r chi.Router) {
	handler := s.newHandler()
	strictHandlers := gen.NewStrictHandlerWithOptions(
		handler,
		[]gen.StrictMiddlewareFunc{},
		gen.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  s.handleRequestError,
			ResponseErrorHandlerFunc: s.handleResponseError,
		},
	)

	gen.HandlerWithOptions(strictHandlers, gen.ChiServerOptions{
		BaseURL:     "/api/v1",
		BaseRouter:  r,
		Middlewares: []gen.MiddlewareFunc{},
	})
}

var _ gen.StrictServerInterface = (*handler)(nil)

type handler struct {
	*errorCodeHandler
	*configHandler
	*systemHandler
	*dashboardDataHandler
	*peripheralHandler
	*commandHandler
}

func (s *Service) newHandler() *handler {
	return &handler{
		errorCodeHandler:     newErrorCodeHandler(s.apperrorcodeService),
		configHandler:        newConfigHandler(s.configService),
		systemHandler:        newSystemHandler(s.systemService),
		dashboardDataHandler: newDashboardDataHandler(s.dashboardDataService),
		peripheralHandler:    newPeripheralHandler(s.peripheralService),
		commandHandler:       newCommandHandler(s.commandService),
	}
}
