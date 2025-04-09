package rfidusb

import (
	"context"
	"log/slog"

	"github.com/tbe-team/raybot/internal/services/location"
)

func (s *Service) HandleRFIDTag(ctx context.Context, tag string) {
	s.log.Debug("RFID tag detected", slog.String("tag", tag))
	if err := s.locationService.UpdateLocation(ctx, location.UpdateLocationParams{
		CurrentLocation: tag,
	}); err != nil {
		s.log.Error("failed to update location", slog.Any("error", err))
	}
}
