package event

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"

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
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
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
			return fmt.Errorf("failed to update ESP serial connection: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		if err := s.systemService.SetStatusError(ctx); err != nil {
			return fmt.Errorf("failed to set system status error: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		s.log.Error("failed to update ESP serial connection", slog.Any("error", err))
	}
}
