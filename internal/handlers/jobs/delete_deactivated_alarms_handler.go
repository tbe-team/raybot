package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/services/alarm"
)

const (
	deleteDeactivatedAlarmsInterval  = 24 * time.Hour
	deleteDeactivatedAlarmsThreshold = 7 * 24 * time.Hour
)

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
	stoppedCh := make(chan struct{})

	go h.run(ctx, stoppedCh)

	return func() {
		cancel()
		<-stoppedCh
	}
}

func (h *deleteDeactivatedAlarmsHandler) run(ctx context.Context, stoppedCh chan struct{}) {
	defer close(stoppedCh)

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(deleteDeactivatedAlarmsInterval):
			if err := h.alarmService.DeleteDeactivatedAlarmsByThreshold(ctx,
				alarm.DeleteDeactivatedAlarmsByThresholdParams{
					Threshold: time.Now().Add(-deleteDeactivatedAlarmsThreshold),
				},
			); err != nil {
				h.log.Error("failed to delete deactivated alarms", slog.Any("error", err))
			}
		}
	}
}
