package event

import (
	"context"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/appstate"
	"github.com/tbe-team/raybot/pkg/ptr"
)

func (s *Service) HandlePICSerialConnectedEvent(ctx context.Context, _ events.PICSerialConnectedEvent) {
	if err := s.appStateService.UpdatePICSerialConnection(ctx, appstate.UpdatePICSerialConnectionParams{
		Connected:          true,
		SetConnected:       true,
		LastConnectedAt:    ptr.New(time.Now()),
		SetLastConnectedAt: true,
	}); err != nil {
		s.log.Error("failed to update PIC serial connection", slog.Any("error", err))
	}
}

func (s *Service) HandlePICSerialDisconnectedEvent(ctx context.Context, event events.PICSerialDisconnectedEvent) {
	var errStr string
	if event.Error != nil {
		errStr = event.Error.Error()
	}

	if err := s.appStateService.UpdatePICSerialConnection(ctx, appstate.UpdatePICSerialConnectionParams{
		Connected:    false,
		SetConnected: true,
		Error:        &errStr,
		SetError:     true,
	}); err != nil {
		s.log.Error("failed to update PIC serial connection", slog.Any("error", err))
	}
}
