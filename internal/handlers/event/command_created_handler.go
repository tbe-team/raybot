package event

import (
	"context"
	"log/slog"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
)

func (s *Service) HandleCommandCreatedEvent(ctx context.Context, event events.CommandCreatedEvent) {
	s.log.Debug("command created event received", slog.Any("event", event))

	if err := s.commandService.ExecuteCreatedCommand(ctx, command.ExecuteCreatedCommandParams{
		CommandID: event.CommandID,
	}); err != nil {
		s.log.Error("failed to execute command",
			slog.Any("event", event),
			slog.Any("error", err),
		)
	}
}
