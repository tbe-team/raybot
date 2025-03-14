package rfid

import (
	"context"
	"log/slog"

	"github.com/tbe-team/raybot/internal/service"
)

type rfidTagHandler struct {
	robotService service.RobotService
	log          *slog.Logger
}

func newRFIDTagHandler(robotService service.RobotService, log *slog.Logger) *rfidTagHandler {
	return &rfidTagHandler{robotService: robotService, log: log}
}

func (h rfidTagHandler) Handle(ctx context.Context, tag string) {
	h.log.Debug("RFID tag detected", slog.String("tag", tag))
	params := service.UpdateRobotStateParams{
		Location: service.LocationParams{
			CurrentLocation: tag,
		},
		SetLocation: true,
	}
	if _, err := h.robotService.UpdateRobotState(ctx, params); err != nil {
		h.log.Error("failed to update robot state", "error", err)
	}
}
