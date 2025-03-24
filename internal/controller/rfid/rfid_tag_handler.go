package rfid

import (
	"context"
	"log/slog"

	"github.com/tbe-team/raybot/internal/service"
)

type rfidTagHandler struct {
	locationService service.LocationService
	log             *slog.Logger
}

func newRFIDTagHandler(locationService service.LocationService, log *slog.Logger) *rfidTagHandler {
	return &rfidTagHandler{locationService: locationService, log: log}
}

func (h rfidTagHandler) Handle(ctx context.Context, tag string) {
	h.log.Debug("RFID tag detected", slog.String("tag", tag))
	params := service.UpdateLocationParams{
		CurrentLocation: tag,
	}
	if err := h.locationService.UpdateLocation(ctx, params); err != nil {
		h.log.Error("failed to update location", "error", err)
	}
}
