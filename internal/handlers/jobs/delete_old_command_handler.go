package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/services/command"
)

type deleteOldCommandHandler struct {
	deleteOldCommandCfg config.DeleteOldCommand

	log            *slog.Logger
	commandService command.Service
}

func newDeleteOldCommandHandler(
	deleteOldCommandCfg config.DeleteOldCommand,
	log *slog.Logger,
	commandService command.Service,
) *deleteOldCommandHandler {
	return &deleteOldCommandHandler{
		deleteOldCommandCfg: deleteOldCommandCfg,
		log:                 log,
		commandService:      commandService,
	}
}

func (h *deleteOldCommandHandler) Run(ctx context.Context) func() {
	ctx, cancel := context.WithCancel(ctx)
	stoppedCh := make(chan struct{})

	go h.run(ctx, stoppedCh)

	return func() {
		cancel()
		<-stoppedCh
	}
}

func (h *deleteOldCommandHandler) run(ctx context.Context, stoppedCh chan struct{}) {
	defer close(stoppedCh)

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(h.deleteOldCommandCfg.ScheduleDuration()):
			if err := h.commandService.DeleteOldCommands(ctx); err != nil {
				h.log.Error("failed to delete old commands", slog.Any("error", err))
			}
		}
	}
}
