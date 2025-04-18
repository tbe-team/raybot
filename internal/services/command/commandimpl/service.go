package commandimpl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/appstate"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/command/executor"
	"github.com/tbe-team/raybot/pkg/eventbus"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/ptr"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	deleteOldCmdCfg config.DeleteOldCommand

	log       *slog.Logger
	validator validator.Validator

	publisher eventbus.Publisher

	runningCmdRepository *runningCmdRepository
	commandRepository    command.Repository
	appStateRepo         appstate.Repository
	executorRouter       executor.Router
}

func NewService(
	deleteOldCmdCfg config.DeleteOldCommand,
	log *slog.Logger,
	validator validator.Validator,
	publisher eventbus.Publisher,
	commandRepository command.Repository,
	appStateRepo appstate.Repository,
	executorRouter executor.Router,
) command.Service {
	s := &service{
		deleteOldCmdCfg:      deleteOldCmdCfg,
		log:                  log.With("service", "command"),
		validator:            validator,
		publisher:            publisher,
		runningCmdRepository: newRunningCmdRepository(),
		commandRepository:    commandRepository,
		appStateRepo:         appStateRepo,
		executorRouter:       executorRouter,
	}

	go s.startCheckingForExecutableCommand(context.Background())

	return s
}

func (s *service) GetCommandByID(ctx context.Context, params command.GetCommandByIDParams) (command.Command, error) {
	if err := s.validator.Validate(params); err != nil {
		return command.Command{}, fmt.Errorf("validate params: %w", err)
	}

	return s.commandRepository.GetCommandByID(ctx, params.CommandID)
}

func (s *service) GetCurrentProcessingCommand(ctx context.Context) (command.Command, error) {
	return s.commandRepository.GetCurrentProcessingCommand(ctx)
}

func (s *service) ListCommands(ctx context.Context, params command.ListCommandsParams) (paging.List[command.Command], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[command.Command]{}, fmt.Errorf("validate params: %w", err)
	}

	return s.commandRepository.ListCommands(ctx, params)
}

func (s *service) CreateCommand(ctx context.Context, params command.CreateCommandParams) (command.Command, error) {
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

func (s *service) CancelCurrentProcessingCommand(_ context.Context) error {
	runningCmd := s.runningCmdRepository.Get()
	if runningCmd == nil {
		return command.ErrNoCommandBeingProcessed
	}

	runningCmd.Cancel()

	return nil
}

func (s *service) ExecuteCreatedCommand(ctx context.Context, params command.ExecuteCreatedCommandParams) error {
	cmd, err := s.commandRepository.GetCommandByID(ctx, params.CommandID)
	if err != nil {
		return fmt.Errorf("get command by id: %w", err)
	}

	if cmd.Status != command.StatusQueued {
		return fmt.Errorf("command is not in queued status")
	}

	processingExists, err := s.commandRepository.CommandProcessingExists(ctx)
	if err != nil {
		return fmt.Errorf("check if command processing exists: %w", err)
	}
	if processingExists {
		s.log.Info("command in PROCESSING status already exists, this command will be queued")
		return nil
	}

	cmd, err = s.commandRepository.UpdateCommand(ctx, command.UpdateCommandParams{
		ID:           cmd.ID,
		Status:       command.StatusProcessing,
		SetStatus:    true,
		StartedAt:    ptr.New(time.Now()),
		SetStartedAt: true,
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		return fmt.Errorf("update command: %w", err)
	}

	s.executeCommand(ctx, cmd)

	return nil
}

func (s *service) DeleteCommandByID(ctx context.Context, params command.DeleteCommandByIDParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	return s.commandRepository.DeleteCommandByIDAndNotProcessing(ctx, params.CommandID)
}

func (s *service) DeleteOldCommands(ctx context.Context) error {
	cutoffTime := time.Now().Add(-s.deleteOldCmdCfg.Threshold)
	return s.commandRepository.DeleteOldCommands(ctx, cutoffTime)
}

func (s *service) startCheckingForExecutableCommand(ctx context.Context) {
	s.log.Debug("waiting for hardware components to be initialized then resuming executable command")
	s.waitForHardwareComponentsInitialized(ctx)

	go s.runNextExecutableCommand(ctx)
}

func (s *service) runNextExecutableCommand(ctx context.Context) {
	cmd, err := s.commandRepository.GetNextExecutableCommand(ctx)
	if err != nil {
		if errors.Is(err, command.ErrNoNextExecutableCommand) {
			return
		}
		s.log.Error("failed to get next executable command", slog.Any("error", err))
		return
	}

	if cmd.Status == command.StatusQueued {
		_, err = s.commandRepository.UpdateCommand(ctx, command.UpdateCommandParams{
			ID:           cmd.ID,
			Status:       command.StatusProcessing,
			SetStatus:    true,
			StartedAt:    ptr.New(time.Now()),
			SetStartedAt: true,
			UpdatedAt:    time.Now(),
		})
		if err != nil {
			s.log.Error("failed to update command to PROCESSING status", slog.Any("error", err))
			return
		}
	}

	s.log.Info("found executable command, executing", slog.Any("command", cmd))
	s.executeCommand(ctx, cmd)
}

func (s *service) waitForHardwareComponentsInitialized(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := s.appStateRepo.ListenForAppStateChanges(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case appState := <-ch:
			if appState.ESPSerialConnection.ServiceInitialized() &&
				appState.PICSerialConnection.ServiceInitialized() &&
				appState.RFIDUSBConnection.ServiceInitialized() {
				return
			}
		}
	}
}

func (s *service) executeCommand(ctx context.Context, cmd command.Command) {
	runningCmd := newRunningCommand(cmd)
	s.runningCmdRepository.Add(runningCmd)
	defer s.runningCmdRepository.Remove()

	// pass the context of the running command to the executor router
	if err := s.executorRouter.Route(runningCmd.Context(), cmd); err != nil {
		s.log.Error("failed to execute command", slog.Any("command", cmd), slog.Any("error", err))
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		s.runNextExecutableCommand(ctx)
	}()
}
