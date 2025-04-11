package commandimpl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

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
	log       *slog.Logger
	validator validator.Validator

	publisher eventbus.Publisher

	commandRepository command.Repository
	appStateRepo      appstate.Repository
	dispatcher        executor.Dispatcher
}

func NewService(
	log *slog.Logger,
	validator validator.Validator,
	publisher eventbus.Publisher,
	commandRepository command.Repository,
	appStateRepo appstate.Repository,
	dispatcher executor.Dispatcher,
) command.Service {
	s := &service{
		log:               log.With("service", "command"),
		validator:         validator,
		publisher:         publisher,
		commandRepository: commandRepository,
		appStateRepo:      appStateRepo,
		dispatcher:        dispatcher,
	}

	go s.startCheckingForExecutableCommand(context.Background())

	return s
}

func (s service) GetCommandByID(ctx context.Context, params command.GetCommandByIDParams) (command.Command, error) {
	if err := s.validator.Validate(params); err != nil {
		return command.Command{}, fmt.Errorf("validate params: %w", err)
	}

	return s.commandRepository.GetCommandByID(ctx, params.CommandID)
}

func (s service) ListCommands(ctx context.Context, params command.ListCommandsParams) (paging.List[command.Command], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[command.Command]{}, fmt.Errorf("validate params: %w", err)
	}

	return s.commandRepository.ListCommands(ctx, params)
}

func (s service) CreateCommand(ctx context.Context, params command.CreateCommandParams) (command.Command, error) {
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

func (s service) ExecuteCreatedCommand(ctx context.Context, params command.ExecuteCreatedCommandParams) error {
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
		ID:        cmd.ID,
		Status:    command.StatusProcessing,
		SetStatus: true,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("update command: %w", err)
	}

	s.executeCommand(ctx, cmd)

	return nil
}

func (s service) startCheckingForExecutableCommand(ctx context.Context) {
	s.log.Debug("waiting for hardware components to be initialized then resuming executable command")
	s.waitForHardwareComponentsInitialized(ctx)

	go s.runNextExecutableCommand(ctx)
}

func (s service) runNextExecutableCommand(ctx context.Context) {
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
			ID:        cmd.ID,
			Status:    command.StatusProcessing,
			SetStatus: true,
			UpdatedAt: time.Now(),
		})
		if err != nil {
			s.log.Error("failed to update command to PROCESSING status", slog.Any("error", err))
			return
		}
	}

	s.log.Info("found executable command, executing", slog.Any("command", cmd))
	s.executeCommand(ctx, cmd)
}

func (s service) waitForHardwareComponentsInitialized(ctx context.Context) {
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

func (s service) executeCommand(ctx context.Context, cmd command.Command) {
	if err := s.dispatcher.Dispatch(ctx, cmd); err != nil {
		_, err = s.commandRepository.UpdateCommand(ctx, command.UpdateCommandParams{
			ID:             cmd.ID,
			Status:         command.StatusFailed,
			SetStatus:      true,
			Error:          ptr.New(err.Error()),
			SetError:       true,
			CompletedAt:    ptr.New(time.Now()),
			SetCompletedAt: true,
			UpdatedAt:      time.Now(),
		})
		if err != nil {
			s.log.Error("failed to update command status to failed", slog.Any("error", err))
		}

	} else {
		_, err := s.commandRepository.UpdateCommand(ctx, command.UpdateCommandParams{
			ID:             cmd.ID,
			Status:         command.StatusSucceeded,
			SetStatus:      true,
			CompletedAt:    ptr.New(time.Now()),
			SetCompletedAt: true,
			UpdatedAt:      time.Now(),
		})
		if err != nil {
			s.log.Error("failed to update command status to succeeded", slog.Any("error", err))
		}
	}

	go s.runNextExecutableCommand(ctx)
}
