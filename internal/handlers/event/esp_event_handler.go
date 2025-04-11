package event

import (
	"context"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/appstate"
	"github.com/tbe-team/raybot/pkg/ptr"
)

func (s *Service) HandleESPSerialConnectedEvent(ctx context.Context, _ events.ESPSerialConnectedEvent) {
	if err := s.appStateService.UpdateESPSerialConnection(ctx, appstate.UpdateESPSerialConnectionParams{
		Connected:          true,
		SetConnected:       true,
		LastConnectedAt:    ptr.New(time.Now()),
		SetLastConnectedAt: true,
	}); err != nil {
		s.log.Error("failed to update ESP serial connection", slog.Any("error", err))
	}
}

func (s *Service) HandleESPSerialDisconnectedEvent(ctx context.Context, event events.ESPSerialDisconnectedEvent) {
	var errStr string
	if event.Error != nil {
		errStr = event.Error.Error()
	}

	if err := s.appStateService.UpdateESPSerialConnection(ctx, appstate.UpdateESPSerialConnectionParams{
		Connected:    false,
		SetConnected: true,
		Error:        &errStr,
		SetError:     true,
	}); err != nil {
		s.log.Error("failed to update ESP serial connection", slog.Any("error", err))
	}
}
