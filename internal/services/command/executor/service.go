package executor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
	"github.com/tbe-team/raybot/pkg/ptr"
)

type Executor[I command.Inputs, O command.Outputs] interface {
	Execute(ctx context.Context, inputs I) (O, error)
}

type Cancelable interface {
	OnCancel(context.Context) error
}

// CommandExecutor is an executor for a command that supports cancellation.
// OnCancel is invoked if the command is canceled during execution.
type CommandExecutor[I command.Inputs, O command.Outputs] interface {
	Executor[I, O]
	Cancelable
}

type service struct {
	log                      *slog.Logger
	runningCommandRepository command.RunningCommandRepository
	commandRepository        command.Repository

	stopMovementExecutor CommandExecutor[command.StopMovementInputs, command.StopMovementOutputs]
	moveToExecutor       CommandExecutor[command.MoveToInputs, command.MoveToOutputs]
	moveForwardExecutor  CommandExecutor[command.MoveForwardInputs, command.MoveForwardOutputs]
	moveBackwardExecutor CommandExecutor[command.MoveBackwardInputs, command.MoveBackwardOutputs]

	cargoOpenExecutor    CommandExecutor[command.CargoOpenInputs, command.CargoOpenOutputs]
	cargoCloseExecutor   CommandExecutor[command.CargoCloseInputs, command.CargoCloseOutputs]
	cargoLiftExecutor    CommandExecutor[command.CargoLiftInputs, command.CargoLiftOutputs]
	cargoLowerExecutor   CommandExecutor[command.CargoLowerInputs, command.CargoLowerOutputs]
	cargoCheckQRExecutor CommandExecutor[command.CargoCheckQRInputs, command.CargoCheckQROutputs]

	scanLocationExecutor CommandExecutor[command.ScanLocationInputs, command.ScanLocationOutputs]
	waitExecutor         CommandExecutor[command.WaitInputs, command.WaitOutputs]

	cancelableMap map[command.CommandType]Cancelable
}

func NewService(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	driveMotorService drivemotor.Service,
	liftMotorService liftmotor.Service,
	cargoService cargo.Service,
	runningCommandRepository command.RunningCommandRepository,
	commandRepository command.Repository,
) command.ExecutorService {
	stopMovementExecutor := newStopMovementExecutor(driveMotorService)
	moveBackwardExecutor := newMoveBackwardExecutor(driveMotorService)
	moveForwardExecutor := newMoveForwardExecutor(driveMotorService)
	moveToExecutor := newMoveToExecutor(log, subscriber, driveMotorService)

	cargoOpenExecutor := newCargoOpenExecutor(log, subscriber, cargoService)
	cargoCloseExecutor := newCargoCloseExecutor(log, subscriber, cargoService)
	cargoLiftExecutor := newCargoLiftExecutor(log, subscriber, liftMotorService)
	cargoLowerExecutor := newCargoLowerExecutor(log, subscriber, liftMotorService)
	cargoCheckQRExecutor := newCargoCheckQRExecutor(log, subscriber)

	scanLocationExecutor := newScanLocationExecutor(log, subscriber, driveMotorService)
	waitExecutor := newWaitExecutor()

	return &service{
		log:                      log,
		runningCommandRepository: runningCommandRepository,
		commandRepository:        commandRepository,

		stopMovementExecutor: stopMovementExecutor,
		moveBackwardExecutor: moveBackwardExecutor,
		moveForwardExecutor:  moveForwardExecutor,
		moveToExecutor:       moveToExecutor,

		cargoOpenExecutor:    cargoOpenExecutor,
		cargoCloseExecutor:   cargoCloseExecutor,
		cargoLiftExecutor:    cargoLiftExecutor,
		cargoLowerExecutor:   cargoLowerExecutor,
		cargoCheckQRExecutor: cargoCheckQRExecutor,

		scanLocationExecutor: scanLocationExecutor,
		waitExecutor:         waitExecutor,

		cancelableMap: map[command.CommandType]Cancelable{
			command.CommandTypeStopMovement: stopMovementExecutor,
			command.CommandTypeMoveBackward: moveBackwardExecutor,
			command.CommandTypeMoveForward:  moveForwardExecutor,
			command.CommandTypeMoveTo:       moveToExecutor,

			command.CommandTypeCargoOpen:    cargoOpenExecutor,
			command.CommandTypeCargoClose:   cargoCloseExecutor,
			command.CommandTypeCargoLift:    cargoLiftExecutor,
			command.CommandTypeCargoLower:   cargoLowerExecutor,
			command.CommandTypeCargoCheckQR: cargoCheckQRExecutor,

			command.CommandTypeScanLocation: scanLocationExecutor,
			command.CommandTypeWait:         waitExecutor,
		},
	}
}

func (s *service) Execute(ctx context.Context, cmd command.Command) error {
	outputs, err := s.execute(ctx, cmd)
	switch {
	case err == nil:
		return s.handleSuccess(ctx, cmd.ID, outputs)

	case errors.Is(err, context.Canceled):
		return s.handleCancel(ctx, cmd.ID, outputs)

	default:
		return s.handleFailure(ctx, cmd.ID, err)
	}
}

func (s *service) execute(ctx context.Context, cmd command.Command) (command.Outputs, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cmd, err := s.commandRepository.UpdateCommand(ctx, command.UpdateCommandParams{
		ID:           cmd.ID,
		Status:       command.StatusProcessing,
		SetStatus:    true,
		StartedAt:    ptr.New(time.Now()),
		SetStartedAt: true,
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update command status: %w", err)
	}

	runningCmd := command.NewCancelableCommand(ctx, cmd)
	defer func() {
		if err := s.runningCommandRepository.Remove(ctx); err != nil {
			s.log.Error("failed to remove running command", slog.Any("error", err))
		}
	}()

	if err := s.runningCommandRepository.Add(ctx, runningCmd); err != nil {
		return nil, fmt.Errorf("failed to add running command: %w", err)
	}

	cmdCtx := runningCmd.Context()
	out, err := s.route(cmdCtx, cmd)

	select {
	case <-cmdCtx.Done():
		s.runCancelHook(ctx, cmd)
		return out, cmdCtx.Err()

	default:
		return out, err
	}
}

func (s *service) runCancelHook(ctx context.Context, cmd command.Command) {
	c, ok := s.cancelableMap[cmd.Type]
	if !ok {
		s.log.Error("cancelable executor not found", slog.Any("command", cmd))
		return
	}
	if err := c.OnCancel(ctx); err != nil {
		s.log.Error("failed to cancel command", slog.Any("command", cmd), slog.Any("error", err))
	}
}

func (s *service) handleSuccess(ctx context.Context, id int64, outputs command.Outputs) error {
	log := s.log.With(slog.Int64("command_id", id), slog.Any("outputs", outputs))
	log.Info("command executed successfully")

	now := time.Now()
	_, err := s.commandRepository.UpdateCommand(ctx, command.UpdateCommandParams{
		ID:             id,
		Status:         command.StatusSucceeded,
		SetStatus:      true,
		Outputs:        outputs,
		SetOutputs:     true,
		CompletedAt:    ptr.New(now),
		SetCompletedAt: true,
		UpdatedAt:      now,
	})

	if err != nil {
		return fmt.Errorf("failed to update command status: %w", err)
	}

	return nil
}

func (s *service) handleCancel(ctx context.Context, id int64, outputs command.Outputs) error {
	log := s.log.With(slog.Int64("command_id", id))
	log.Info("command cancelled")

	now := time.Now()
	_, err := s.commandRepository.UpdateCommand(ctx, command.UpdateCommandParams{
		ID:             id,
		Status:         command.StatusCanceled,
		SetStatus:      true,
		Outputs:        outputs,
		SetOutputs:     true,
		CompletedAt:    ptr.New(now),
		SetCompletedAt: true,
		UpdatedAt:      now,
	})
	if err != nil {
		return fmt.Errorf("failed to update command status: %w", err)
	}

	return nil
}

func (s *service) handleFailure(ctx context.Context, id int64, execErr error) error {
	log := s.log.With(slog.Int64("command_id", id), slog.Any("exec_error", execErr))
	log.Error("command execution failed")

	now := time.Now()
	_, err := s.commandRepository.UpdateCommand(ctx, command.UpdateCommandParams{
		ID:             id,
		Status:         command.StatusFailed,
		SetStatus:      true,
		Error:          ptr.New(execErr.Error()),
		SetError:       true,
		CompletedAt:    ptr.New(now),
		SetCompletedAt: true,
		UpdatedAt:      now,
	})
	if err != nil {
		return fmt.Errorf("failed to update command status: %w", err)
	}

	return nil
}
