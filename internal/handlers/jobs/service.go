package jobs

import (
	"context"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/services/alarm"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type Service struct {
	cronCfg config.Cron
	log     *slog.Logger

	subscriber     eventbus.Subscriber
	commandService command.Service
	alarmService   alarm.Service
}

type CleanupFunc func(context.Context) error

func New(
	cronCfg config.Cron,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	commandService command.Service,
	alarmService alarm.Service,
) *Service {
	return &Service{
		cronCfg:        cronCfg,
		log:            log.With("service", "jobs"),
		subscriber:     subscriber,
		commandService: commandService,
		alarmService:   alarmService,
	}
}

func (s *Service) Run(ctx context.Context) (CleanupFunc, error) {
	deleteOldCommandHandler := newDeleteOldCommandHandler(s.cronCfg.DeleteOldCommand, s.log, s.commandService)
	executeCommandHandler := newExecuteCommandHandler(s.log, s.commandService, s.subscriber)
	deleteDeactivatedAlarmsHandler := newDeleteDeactivatedAlarmsHandler(s.log, s.alarmService)

	stopFuncs := []func(){}
	stopFuncs = append(stopFuncs, deleteOldCommandHandler.Run(ctx))
	stopFuncs = append(stopFuncs, executeCommandHandler.Run(ctx))
	stopFuncs = append(stopFuncs, deleteDeactivatedAlarmsHandler.Run(ctx))

	cleanup := func(_ context.Context) error {
		wg := sync.WaitGroup{}
		for _, stopFunc := range stopFuncs {
			wg.Add(1)
			go func(stopFunc func()) {
				defer wg.Done()
				stopFunc()
			}(stopFunc)
		}

		wg.Wait()

		return nil
	}

	return cleanup, nil
}
