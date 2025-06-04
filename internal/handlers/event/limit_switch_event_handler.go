package event

import (
	"context"
	"log/slog"

	"github.com/tbe-team/raybot/internal/events"
)

func (s *Service) HandleLimitSwitch1PressedEvent(ctx context.Context, _ events.LimitSwitch1PressedEvent) {
	if err := s.commandService.CancelCurrentProcessingCommand(ctx); err != nil {
		s.log.Error("failed to cancel current processing command", slog.Any("error", err))
	}
}
