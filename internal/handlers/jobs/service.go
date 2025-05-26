package jobs

import (
	"context"
	"log/slog"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type Service struct {
	cronCfg config.Cron
	log     *slog.Logger

	subscriber     eventbus.Subscriber
	commandService command.Service
}

type CleanupFunc func(context.Context) error

func New(
	cronCfg config.Cron,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	commandService command.Service,
) *Service {
	return &Service{
		cronCfg:        cronCfg,
		log:            log.With("service", "jobs"),
		subscriber:     subscriber,
		commandService: commandService,
	}
}

func (s *Service) Run(ctx context.Context) (CleanupFunc, error) {
	deleteOldCommandHandler := newDeleteOldCommandHandler(s.cronCfg.DeleteOldCommand, s.log, s.commandService)
	executeCommandHandler := newExecuteCommandHandler(s.log, s.commandService, s.subscriber)

	cancelDeleteOldCommand := deleteOldCommandHandler.Run(ctx)
	cancelExecuteCommand := executeCommandHandler.Run(ctx)

	cleanup := func(_ context.Context) error {
		cancelDeleteOldCommand()
		cancelExecuteCommand()

		return nil
	}

	return cleanup, nil
}
