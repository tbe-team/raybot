package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/services/alarm"
)

const deleteDeactivatedAlarmsInterval = 24 * time.Hour

type deleteDeactivatedAlarmsHandler struct {
	log          *slog.Logger
	alarmService alarm.Service
}

func newDeleteDeactivatedAlarmsHandler(
	log *slog.Logger,
	alarmService alarm.Service,
) *deleteDeactivatedAlarmsHandler {
	return &deleteDeactivatedAlarmsHandler{
		log:          log,
		alarmService: alarmService,
	}
}

func (h *deleteDeactivatedAlarmsHandler) Run(ctx context.Context) func() {
	ctx, cancel := context.WithCancel(ctx)

	go h.run(ctx)

	return cancel
}

func (h *deleteDeactivatedAlarmsHandler) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(deleteDeactivatedAlarmsInterval):
			if err := h.alarmService.DeleteDeactivatedAlarms(ctx); err != nil {
				h.log.Error("failed to delete deactivated alarms", slog.Any("error", err))
			}
		}
	}
}
