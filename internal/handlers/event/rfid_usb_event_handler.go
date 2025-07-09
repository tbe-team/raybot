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

func (s *Service) HandleRFIDUSBConnectedEvent(ctx context.Context, _ events.RFIDUSBConnectedEvent) {
	if err := s.appStateService.UpdateRFIDUSBConnection(ctx, appstate.UpdateRFIDUSBConnectionParams{
		Connected:          true,
		SetConnected:       true,
		LastConnectedAt:    ptr.New(time.Now()),
		SetLastConnectedAt: true,
	}); err != nil {
		s.log.Error("failed to update RFID USB connection", slog.Any("error", err))
	}
}

func (s *Service) HandleRFIDUSBDisconnectedEvent(ctx context.Context, event events.RFIDUSBDisconnectedEvent) {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var errStr string
		if event.Error != nil {
			errStr = event.Error.Error()
		}

		if err := s.appStateService.UpdateRFIDUSBConnection(ctx, appstate.UpdateRFIDUSBConnectionParams{
			Connected:    false,
			SetConnected: true,
			Error:        &errStr,
			SetError:     true,
		}); err != nil {
			return fmt.Errorf("failed to update RFID USB connection: %w", err)
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
		s.log.Error("failed to update cloud connection", slog.Any("error", err))
	}
}
