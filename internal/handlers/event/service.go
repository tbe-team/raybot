package event

import (
	"context"
	"log/slog"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/appstate"
)

type Service struct {
	log *slog.Logger

	appStateService appstate.Service
}

type CleanupFunc func(context.Context) error

func New(
	log *slog.Logger,
	appStateService appstate.Service,
) *Service {
	return &Service{
		log:             log.With("service", "event"),
		appStateService: appStateService,
	}
}

func (s *Service) Run(_ context.Context) (CleanupFunc, error) {
	events.CloudConnectedSignal.AddListener(s.HandleCloudConnectedEvent)
	events.CloudDisconnectedSignal.AddListener(s.HandleCloudDisconnectedEvent)
	events.ESPSerialConnectedSignal.AddListener(s.HandleESPSerialConnectedEvent)
	events.ESPSerialDisconnectedSignal.AddListener(s.HandleESPSerialDisconnectedEvent)
	events.PICSerialConnectedSignal.AddListener(s.HandlePICSerialConnectedEvent)
	events.PICSerialDisconnectedSignal.AddListener(s.HandlePICSerialDisconnectedEvent)
	events.RFIDUSBConnectedSignal.AddListener(s.HandleRFIDUSBConnectedEvent)
	events.RFIDUSBDisconnectedSignal.AddListener(s.HandleRFIDUSBDisconnectedEvent)

	cleanup := func(_ context.Context) error {
		return nil
	}

	return cleanup, nil
}
