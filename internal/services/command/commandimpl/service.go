package commandimpl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/pkg/eventbus"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/validator"
)

type Service struct {
	deleteOldCmdCfg config.DeleteOldCommand

	log       *slog.Logger
	validator validator.Validator

	publisher eventbus.Publisher

	runningCmdRepository command.RunningCommandRepository
	commandRepository    command.Repository

	processingLock  command.ProcessingLock
	executorService command.ExecutorService
}

func NewService(
	deleteOldCmdCfg config.DeleteOldCommand,
	log *slog.Logger,
	validator validator.Validator,
	publisher eventbus.Publisher,
	runningCmdRepository command.RunningCommandRepository,
	commandRepository command.Repository,
	processingLock command.ProcessingLock,
	executorService command.ExecutorService,
) command.Service {
	s := &Service{
		deleteOldCmdCfg:      deleteOldCmdCfg,
		log:                  log.With("service", "command"),
		validator:            validator,
		publisher:            publisher,
		runningCmdRepository: runningCmdRepository,
		commandRepository:    commandRepository,
		processingLock:       processingLock,
		executorService:      executorService,
	}

	go s.cancelQueuedAndProcessingCommands(context.Background())

	return s
}

func (s *Service) GetCommandByID(ctx context.Context, params command.GetCommandByIDParams) (command.Command, error) {
	if err := s.validator.Validate(params); err != nil {
		return command.Command{}, fmt.Errorf("validate params: %w", err)
	}

	return s.commandRepository.GetCommandByID(ctx, params.CommandID)
}

func (s *Service) GetCurrentProcessingCommand(ctx context.Context) (command.Command, error) {
	return s.commandRepository.GetCurrentProcessingCommand(ctx)
}

func (s *Service) ListCommands(ctx context.Context, params command.ListCommandsParams) (paging.List[command.Command], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[command.Command]{}, fmt.Errorf("validate params: %w", err)
	}

	return s.commandRepository.ListCommands(ctx, params)
}

func (s *Service) CreateCommand(ctx context.Context, params command.CreateCommandParams) (command.Command, error) {
	if err := s.validator.Validate(params); err != nil {
		return command.Command{}, fmt.Errorf("validate params: %w", err)
	}

	cmd := command.NewCommand(params.Source, params.Inputs)
	cmd, err := s.commandRepository.CreateCommand(ctx, cmd)
	if err != nil {
		return command.Command{}, fmt.Errorf("create command: %w", err)
	}

	s.publisher.Publish(
		events.CommandCreatedTopic,
		eventbus.NewMessage(events.CommandCreatedEvent{
			CommandID: cmd.ID,
		}),
	)

	return cmd, nil
}

func (s *Service) CancelCurrentProcessingCommand(ctx context.Context) error {
	runningCmd, err := s.runningCmdRepository.Get(ctx)
	if err != nil {
		if errors.Is(err, command.ErrRunningCommandNotFound) {
			return command.ErrNoCommandBeingProcessed
		}
		return fmt.Errorf("get running command: %w", err)
	}

	runningCmd.Cancel()

	return nil
}

func (s *Service) CancelActiveCloudCommands(ctx context.Context) error {
	if err := s.processingLock.WithLock(func() error {
		// Cancel current processing command
		runningCmd, err := s.runningCmdRepository.Get(ctx)
		if err != nil {
			if !errors.Is(err, command.ErrRunningCommandNotFound) {
				return fmt.Errorf("get running command: %w", err)
			}
		} else {
			if runningCmd.Source == command.SourceCloud {
				runningCmd.Cancel()
			}
		}

		// Cancel all queued and processing commands created by the cloud
		if err := s.commandRepository.CancelQueuedAndProcessingCommandsCreatedByCloud(ctx); err != nil {
			return fmt.Errorf("cancel queued and processing commands created by cloud: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("cancel active cloud commands: %w", err)
	}

	return nil
}

func (s *Service) ExecuteCreatedCommand(ctx context.Context, params command.ExecuteCreatedCommandParams) error {
	_, err := s.runningCmdRepository.Get(ctx)
	// no error means the running command exists, so we don't need to execute the command
	if err == nil {
		s.log.Info("command is already being processed, this command will be queued")
		return nil
	}

	cmd, err := s.commandRepository.GetCommandByID(ctx, params.CommandID)
	if err != nil {
		return fmt.Errorf("get command by id: %w", err)
	}

	if cmd.Status != command.StatusQueued {
		return fmt.Errorf("command is not in queued status")
	}

	s.executeCommand(ctx, cmd)

	return nil
}

func (s *Service) CancelAllRunningCommands(ctx context.Context) error {
	if err := s.processingLock.WithLock(func() error {
		runningCmd, err := s.runningCmdRepository.Get(ctx)
		if err != nil {
			if !errors.Is(err, command.ErrRunningCommandNotFound) {
				return fmt.Errorf("get running command: %w", err)
			}
		} else {
			runningCmd.Cancel()
			if err := s.runningCmdRepository.Remove(ctx); err != nil {
				return fmt.Errorf("remove running command: %w", err)
			}
		}

		if err := s.commandRepository.CancelQueuedAndProcessingCommands(ctx); err != nil {
			return fmt.Errorf("cancel queued and processing commands: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("with processing lock: %w", err)
	}

	return nil
}

func (s *Service) DeleteCommandByID(ctx context.Context, params command.DeleteCommandByIDParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.commandRepository.DeleteCommandByIDAndNotProcessing(ctx, params.CommandID)
}

func (s *Service) DeleteOldCommands(ctx context.Context) error {
	cutoffTime := time.Now().Add(-s.deleteOldCmdCfg.Threshold)
	return s.commandRepository.DeleteOldCommands(ctx, cutoffTime)
}

func (s *Service) cancelQueuedAndProcessingCommands(ctx context.Context) {
	if err := s.commandRepository.CancelQueuedAndProcessingCommands(ctx); err != nil {
		s.log.Error("failed to cancel queued and processing commands on startup", slog.Any("error", err))
	}
}

func (s *Service) runNextExecutableCommand(ctx context.Context) {
	cmd, err := s.commandRepository.GetNextExecutableCommand(ctx)
	if err != nil {
		if errors.Is(err, command.ErrNoNextExecutableCommand) {
			return
		}
		s.log.Error("failed to get next executable command", slog.Any("error", err))
		return
	}

	s.log.Info("found executable command, executing", slog.Any("command", cmd))
	s.executeCommand(ctx, cmd)
}

func (s *Service) executeCommand(ctx context.Context, cmd command.Command) {
	if err := s.processingLock.WaitUntilUnlocked(ctx); err != nil {
		s.log.Error("failed to wait for processing lock to be unlocked", slog.Any("error", err))
	}

	if err := s.executorService.Execute(ctx, cmd); err != nil {
		s.log.Error("failed to execute command", slog.Any("command", cmd), slog.Any("error", err))
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		s.runNextExecutableCommand(ctx)
	}()
}
