package jobs

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-co-op/gocron/v2"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/services/command"
)

type Service struct {
	cronCfg config.Cron
	log     *slog.Logger

	commandService command.Service
}

type CleanupFunc func(context.Context) error

func New(
	cronCfg config.Cron,
	log *slog.Logger,
	commandService command.Service,
) *Service {
	return &Service{
		cronCfg:        cronCfg,
		log:            log.With("service", "jobs"),
		commandService: commandService,
	}
}

func (s *Service) Run() (CleanupFunc, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	if err := s.scheduleCronJobs(scheduler); err != nil {
		return nil, fmt.Errorf("failed to schedule cron jobs: %w", err)
	}

	scheduler.Start()

	cleanup := func(_ context.Context) error {
		return scheduler.Shutdown()
	}

	return cleanup, nil
}

func (s *Service) scheduleCronJobs(scheduler gocron.Scheduler) error {
	_, err := scheduler.NewJob(
		gocron.CronJob(s.cronCfg.DeleteOldCommand.Schedule, true),
		gocron.NewTask(func(ctx context.Context) {
			s.log.Debug("running delete old commands job")
			if err := s.HandleDeleteOldCommands(ctx); err != nil {
				s.log.Error("failed to handle delete old commands", slog.Any("error", err))
			}
		}),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)
	if err != nil {
		return fmt.Errorf("failed to schedule delete old commands: %w", err)
	}

	return nil
}
