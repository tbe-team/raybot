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

func (s *Service) HandleCloudConnectedEvent(ctx context.Context, _ events.CloudConnectedEvent) {
	if err := s.appStateService.UpdateCloudConnection(ctx, appstate.UpdateCloudConnectionParams{
		Connected:          true,
		SetConnected:       true,
		LastConnectedAt:    ptr.New(time.Now()),
		SetLastConnectedAt: true,
	}); err != nil {
		s.log.Error("failed to update cloud connection", slog.Any("error", err))
	}
}

func (s *Service) HandleCloudDisconnectedEvent(ctx context.Context, event events.CloudDisconnectedEvent) {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var errStr string
		if event.Error != nil {
			errStr = event.Error.Error()
		}

		if err := s.appStateService.UpdateCloudConnection(ctx, appstate.UpdateCloudConnectionParams{
			Connected:    false,
			SetConnected: true,
			Error:        &errStr,
			SetError:     true,
		}); err != nil {
			return fmt.Errorf("failed to update cloud connection: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		if err := s.commandService.CancelActiveCloudCommands(ctx); err != nil {
			return fmt.Errorf("failed to cancel active cloud commands: %w", err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		s.log.Error("failed to update cloud connection", slog.Any("error", err))
	}
}
