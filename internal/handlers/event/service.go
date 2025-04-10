package event

import (
	"context"
	"log/slog"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/appstate"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type Service struct {
	log *slog.Logger

	subscriber eventbus.Subscriber

	appStateService appstate.Service
	commandService  command.Service
}

type CleanupFunc func(context.Context) error

func New(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	appStateService appstate.Service,
	commandService command.Service,
) *Service {
	return &Service{
		log:             log.With("service", "event"),
		subscriber:      subscriber,
		appStateService: appStateService,
		commandService:  commandService,
	}
}

func (s *Service) Run(ctx context.Context) (CleanupFunc, error) {
	ctx, cancel := context.WithCancel(ctx)
	s.registerHandlers(ctx)

	cleanup := func(_ context.Context) error {
		cancel()
		return nil
	}

	return cleanup, nil
}

func (s *Service) registerHandlers(ctx context.Context) {
	s.subscriber.Subscribe(
		ctx,
		events.CloudConnectedTopic,
		func(ctx context.Context, msg *eventbus.Message) {
			ev, ok := msg.Payload.(events.CloudConnectedEvent)
			if !ok {
				s.log.Error("received invalid event", slog.Any("event", msg.Payload))
				return
			}

			s.HandleCloudConnectedEvent(ctx, ev)
		},
	)

	s.subscriber.Subscribe(
		ctx,
		events.CloudDisconnectedTopic,
		func(ctx context.Context, msg *eventbus.Message) {
			ev, ok := msg.Payload.(events.CloudDisconnectedEvent)
			if !ok {
				s.log.Error("received invalid event", slog.Any("event", msg.Payload))
				return
			}

			s.HandleCloudDisconnectedEvent(ctx, ev)
		},
	)

	s.subscriber.Subscribe(
		ctx,
		events.ESPSerialConnectedTopic,
		func(ctx context.Context, msg *eventbus.Message) {
			ev, ok := msg.Payload.(events.ESPSerialConnectedEvent)
			if !ok {
				s.log.Error("received invalid event", slog.Any("event", msg.Payload))
				return
			}

			s.HandleESPSerialConnectedEvent(ctx, ev)
		},
	)

	s.subscriber.Subscribe(
		ctx,
		events.ESPSerialDisconnectedTopic,
		func(ctx context.Context, msg *eventbus.Message) {
			ev, ok := msg.Payload.(events.ESPSerialDisconnectedEvent)
			if !ok {
				s.log.Error("received invalid event", slog.Any("event", msg.Payload))
				return
			}

			s.HandleESPSerialDisconnectedEvent(ctx, ev)
		},
	)

	s.subscriber.Subscribe(
		ctx,
		events.PICSerialConnectedTopic,
		func(ctx context.Context, msg *eventbus.Message) {
			ev, ok := msg.Payload.(events.PICSerialConnectedEvent)
			if !ok {
				s.log.Error("received invalid event", slog.Any("event", msg.Payload))
				return
			}

			s.HandlePICSerialConnectedEvent(ctx, ev)
		},
	)

	s.subscriber.Subscribe(
		ctx,
		events.PICSerialDisconnectedTopic,
		func(ctx context.Context, msg *eventbus.Message) {
			ev, ok := msg.Payload.(events.PICSerialDisconnectedEvent)
			if !ok {
				s.log.Error("received invalid event", slog.Any("event", msg.Payload))
				return
			}

			s.HandlePICSerialDisconnectedEvent(ctx, ev)
		},
	)

	s.subscriber.Subscribe(
		ctx,
		events.RFIDUSBConnectedTopic,
		func(ctx context.Context, msg *eventbus.Message) {
			ev, ok := msg.Payload.(events.RFIDUSBConnectedEvent)
			if !ok {
				s.log.Error("received invalid event", slog.Any("event", msg.Payload))
				return
			}

			s.HandleRFIDUSBConnectedEvent(ctx, ev)
		},
	)

	s.subscriber.Subscribe(
		ctx,
		events.RFIDUSBDisconnectedTopic,
		func(ctx context.Context, msg *eventbus.Message) {
			ev, ok := msg.Payload.(events.RFIDUSBDisconnectedEvent)
			if !ok {
				s.log.Error("received invalid event", slog.Any("event", msg.Payload))
				return
			}

			s.HandleRFIDUSBDisconnectedEvent(ctx, ev)
		},
	)

	s.subscriber.Subscribe(
		ctx,
		events.CommandCreatedTopic,
		func(ctx context.Context, msg *eventbus.Message) {
			ev, ok := msg.Payload.(events.CommandCreatedEvent)
			if !ok {
				s.log.Error("received invalid event", slog.Any("event", msg.Payload))
				return
			}

			s.HandleCommandCreatedEvent(ctx, ev)
		},
	)
}
