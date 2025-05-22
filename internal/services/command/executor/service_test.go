package executor

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/logging"
	cargomocks "github.com/tbe-team/raybot/internal/services/cargo/mocks"
	"github.com/tbe-team/raybot/internal/services/command"
	commandmocks "github.com/tbe-team/raybot/internal/services/command/mocks"
	drivemotormocks "github.com/tbe-team/raybot/internal/services/drivemotor/mocks"
	liftmotormocks "github.com/tbe-team/raybot/internal/services/liftmotor/mocks"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

func TestService_NewService(t *testing.T) {
	service := NewService(
		logging.NewNoopLogger(),
		&eventbus.NoopEventBus{},
		drivemotormocks.NewFakeService(t),
		liftmotormocks.NewFakeService(t),
		cargomocks.NewFakeService(t),
		commandmocks.NewFakeRunningCommandRepository(t),
		commandmocks.NewFakeRepository(t),
	)
	require.NotNil(t, service)
}

func TestService_Execute(t *testing.T) {
	t.Run("Should execute command successfully", func(t *testing.T) {
		log := logging.NewNoopLogger()
		runningCommandRepository := commandmocks.NewFakeRunningCommandRepository(t)
		commandRepository := commandmocks.NewFakeRepository(t)
		service := newTestService(log, runningCommandRepository, commandRepository, nil)

		cmdID := int64(1)

		commandRepository.EXPECT().UpdateCommand(mock.Anything, mock.MatchedBy(
			func(params command.UpdateCommandParams) bool {
				return params.ID == cmdID &&
					params.Status == command.StatusProcessing &&
					params.SetStatus &&
					params.StartedAt != nil &&
					params.SetStartedAt &&
					!params.UpdatedAt.IsZero()
			},
		)).Return(command.Command{
			ID:     cmdID,                           // Required by the next repository mock call
			Type:   command.CommandTypeStopMovement, // Required by executor
			Inputs: &command.StopMovementInputs{},   // Required by executor
		}, nil)

		runningCommandRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		runningCommandRepository.EXPECT().Remove(mock.Anything).Return(nil)
		commandRepository.EXPECT().UpdateCommand(mock.Anything, mock.MatchedBy(
			func(params command.UpdateCommandParams) bool {
				return params.ID == cmdID &&
					params.Status == command.StatusSucceeded &&
					params.SetStatus &&
					params.SetOutputs &&
					params.Error == nil &&
					!params.SetError &&
					params.CompletedAt != nil &&
					params.SetCompletedAt &&
					!params.UpdatedAt.IsZero()
			},
		)).Return(command.Command{}, nil)

		err := service.Execute(context.Background(), command.Command{
			ID:     cmdID,
			Type:   command.CommandTypeStopMovement,
			Inputs: &command.StopMovementInputs{},
		})
		require.NoError(t, err)
	})

	t.Run("Should fail to execute command and update command status to FAILED successfully", func(t *testing.T) {
		log := logging.NewNoopLogger()
		runningCommandRepository := commandmocks.NewFakeRunningCommandRepository(t)
		commandRepository := commandmocks.NewFakeRepository(t)
		execErr := errors.New("exec error")
		service := newTestService(log, runningCommandRepository, commandRepository, execErr)

		cmdID := int64(1)

		commandRepository.EXPECT().UpdateCommand(mock.Anything, mock.MatchedBy(
			func(params command.UpdateCommandParams) bool {
				return params.ID == cmdID &&
					params.Status == command.StatusProcessing &&
					params.SetStatus &&
					params.StartedAt != nil &&
					params.SetStartedAt &&
					!params.UpdatedAt.IsZero()
			},
		)).Return(command.Command{
			ID:     cmdID,                           // Required by the next repository mock call
			Type:   command.CommandTypeStopMovement, // Required by executor
			Inputs: &command.StopMovementInputs{},   // Required by executor
		}, nil)

		runningCommandRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		runningCommandRepository.EXPECT().Remove(mock.Anything).Return(nil)
		commandRepository.EXPECT().UpdateCommand(mock.Anything, mock.MatchedBy(
			func(params command.UpdateCommandParams) bool {
				return params.ID == cmdID &&
					params.Status == command.StatusFailed &&
					params.SetStatus &&
					!params.SetOutputs &&
					params.Error != nil &&
					*params.Error == execErr.Error() &&
					params.SetError &&
					params.CompletedAt != nil &&
					params.SetCompletedAt &&
					!params.UpdatedAt.IsZero()
			},
		)).Return(command.Command{}, nil)

		err := service.Execute(context.Background(), command.Command{
			ID:     cmdID,
			Type:   command.CommandTypeStopMovement,
			Inputs: &command.StopMovementInputs{},
		})
		require.NoError(t, err)
	})

	t.Run("Should handle command has been canceled and update status to CANCELED successfully", func(t *testing.T) {
		log := logging.NewNoopLogger()
		runningCommandRepository := commandmocks.NewFakeRunningCommandRepository(t)
		commandRepository := commandmocks.NewFakeRepository(t)
		service := newTestService(log, runningCommandRepository, commandRepository, context.Canceled)

		cmdID := int64(1)

		commandRepository.EXPECT().UpdateCommand(mock.Anything, mock.MatchedBy(
			func(params command.UpdateCommandParams) bool {
				return params.ID == cmdID &&
					params.Status == command.StatusProcessing &&
					params.SetStatus &&
					params.StartedAt != nil &&
					params.SetStartedAt &&
					!params.UpdatedAt.IsZero()
			},
		)).Return(command.Command{
			ID:     cmdID,                           // Required by the next repository mock call
			Type:   command.CommandTypeStopMovement, // Required by executor
			Inputs: &command.StopMovementInputs{},   // Required by executor
		}, nil)

		runningCommandRepository.EXPECT().Add(mock.Anything, mock.Anything).Return(nil)
		runningCommandRepository.EXPECT().Remove(mock.Anything).Return(nil)
		commandRepository.EXPECT().UpdateCommand(mock.Anything, mock.MatchedBy(
			func(params command.UpdateCommandParams) bool {
				return params.ID == cmdID &&
					params.Status == command.StatusCanceled &&
					params.SetStatus &&
					params.SetOutputs &&
					params.Outputs != nil &&
					params.CompletedAt != nil &&
					params.SetCompletedAt &&
					!params.UpdatedAt.IsZero()
			},
		)).Return(command.Command{}, nil)

		err := service.Execute(context.Background(), command.Command{
			ID:     cmdID,
			Type:   command.CommandTypeStopMovement,
			Inputs: &command.StopMovementInputs{},
		})
		require.NoError(t, err)
	})
}

func newTestService(
	log *slog.Logger,
	runningCommandRepository command.RunningCommandRepository,
	commandRepository command.Repository,
	expectedReturnErr error,
) *service {
	stopMovementExecutor := newFakeExecutor[command.StopMovementInputs, command.StopMovementOutputs](expectedReturnErr)
	moveBackwardExecutor := newFakeExecutor[command.MoveBackwardInputs, command.MoveBackwardOutputs](expectedReturnErr)
	moveForwardExecutor := newFakeExecutor[command.MoveForwardInputs, command.MoveForwardOutputs](expectedReturnErr)
	moveToExecutor := newFakeExecutor[command.MoveToInputs, command.MoveToOutputs](expectedReturnErr)

	cargoOpenExecutor := newFakeExecutor[command.CargoOpenInputs, command.CargoOpenOutputs](expectedReturnErr)
	cargoCloseExecutor := newFakeExecutor[command.CargoCloseInputs, command.CargoCloseOutputs](expectedReturnErr)
	cargoLiftExecutor := newFakeExecutor[command.CargoLiftInputs, command.CargoLiftOutputs](expectedReturnErr)
	cargoLowerExecutor := newFakeExecutor[command.CargoLowerInputs, command.CargoLowerOutputs](expectedReturnErr)
	cargoCheckQRExecutor := newFakeExecutor[command.CargoCheckQRInputs, command.CargoCheckQROutputs](expectedReturnErr)

	scanLocationExecutor := newFakeExecutor[command.ScanLocationInputs, command.ScanLocationOutputs](expectedReturnErr)
	waitExecutor := newFakeExecutor[command.WaitInputs, command.WaitOutputs](expectedReturnErr)

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

type fakeExecutor[I command.Inputs, O command.Outputs] struct {
	expectedReturnErr error
}

func newFakeExecutor[I command.Inputs, O command.Outputs](returnErr error) CommandExecutor[I, O] {
	return fakeExecutor[I, O]{
		expectedReturnErr: returnErr,
	}
}

func (e fakeExecutor[I, O]) Execute(_ context.Context, _ I) (O, error) {
	var zero O
	return zero, e.expectedReturnErr
}

func (e fakeExecutor[I, O]) OnCancel(_ context.Context) error {
	return nil
}
