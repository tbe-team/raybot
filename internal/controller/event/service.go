package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"

	"github.com/tbe-team/raybot/internal/pubsub"
	"github.com/tbe-team/raybot/internal/service"
)

type CleanupFunc func(ctx context.Context) error

type Service struct {
	service service.Service
	pubsub  pubsub.PubSub
	log     *slog.Logger
}

func New(service service.Service, pubsub pubsub.PubSub, log *slog.Logger) *Service {
	return &Service{
		service: service,
		pubsub:  pubsub,
		log:     log.With(slog.String("service", "event")),
	}
}

func (s Service) Run(ctx context.Context) (CleanupFunc, error) {
	router, err := message.NewRouter(message.RouterConfig{}, watermill.NewSlogLogger(s.log))
	if err != nil {
		return nil, fmt.Errorf("new router: %w", err)
	}

	s.registerRouterMiddleware(router)
	s.registerEventRouter(router)

	go func() {
		s.log.Info("starting event service")
		if err := router.Run(ctx); err != nil {
			s.log.Error("router run", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	cleanup := func(_ context.Context) error {
		s.log.Debug("stopping event service")
		if err := router.Close(); err != nil {
			return fmt.Errorf("close router: %w", err)
		}
		s.log.Debug("event service stopped")
		return nil
	}

	return cleanup, nil
}

func (s Service) registerRouterMiddleware(router *message.Router) {
	router.AddMiddleware(middleware.Recoverer)
}

func (s Service) registerEventRouter(router *message.Router) {
	commandCreatedEventHandler := NewCommandCreatedEventHandler(s.service.CommandService())

	router.AddNoPublisherHandler(
		"command-created-event-handler",
		pubsub.TopicCommandCreated,
		s.pubsub,
		func(msg *message.Message) error {
			var cmdCreatedEvent pubsub.CommandCreatedEvent
			if err := json.Unmarshal(msg.Payload, &cmdCreatedEvent); err != nil {
				return fmt.Errorf("unmarshal command created event: %w", err)
			}

			if err := commandCreatedEventHandler.Handle(msg.Context(), cmdCreatedEvent); err != nil {
				// We don't want to retry the command
				s.log.Error("command created event handler", slog.Any("error", err))
			}

			return nil
		},
	)

}
