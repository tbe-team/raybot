package job

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-co-op/gocron/v2"

	"github.com/tbe-team/raybot/internal/service"
)

type Jobs struct {
	executeInProgressCommand *ExecuteInProgressCommand
}

type Service struct {
	jobs    Jobs
	service service.Service
	log     *slog.Logger
}

func NewService(service service.Service, log *slog.Logger) *Service {
	return &Service{
		jobs: Jobs{
			executeInProgressCommand: NewExecuteInProgressCommand(service.CommandService()),
		},
		service: service,
		log:     log,
	}
}

type CleanupFunc func(ctx context.Context) error

func (s Service) Run() (CleanupFunc, error) {
	s.log.Info("starting scheduler")
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	if _, err := scheduler.NewJob(
		// Run every second
		gocron.CronJob("* * * * * *", true),
		gocron.NewTask(func() {
			if err := s.jobs.executeInProgressCommand.Run(context.Background()); err != nil {
				s.log.Error("failed to execute in progress command", "error", err)
			}
		}),
	); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	scheduler.Start()

	return func(_ context.Context) error {
		s.log.Debug("shutting down scheduler")
		if err := scheduler.Shutdown(); err != nil {
			return fmt.Errorf("failed to shutdown scheduler: %w", err)
		}
		s.log.Debug("scheduler shut down complete")
		return nil
	}, nil
}
