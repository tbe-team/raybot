package commandimpl

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/logging"
	"github.com/tbe-team/raybot/internal/services/command"
	commandmocks "github.com/tbe-team/raybot/internal/services/command/mocks"
	eventbusmocks "github.com/tbe-team/raybot/pkg/eventbus/mocks"
	"github.com/tbe-team/raybot/pkg/validator"
)

func TestService_CreateCommand(t *testing.T) {
	publisher := eventbusmocks.NewFakePublisher(t)
	commandRepository := commandmocks.NewFakeRepository(t)
	commandService := Service{
		validator:         validator.New(),
		publisher:         publisher,
		commandRepository: commandRepository,
	}

	t.Run("Create command successfully", func(t *testing.T) {
		commandRepository.EXPECT().CreateCommand(mock.Anything, mock.Anything).Return(command.Command{}, nil)
		publisher.EXPECT().Publish(events.CommandCreatedTopic, mock.Anything).Once()
		command, err := commandService.CreateCommand(context.Background(), command.CreateCommandParams{
			Source: command.SourceApp,
			Inputs: command.StopMovementInputs{},
		})
		require.NoError(t, err)
		require.NotNil(t, command)
	})

	t.Run("Create command validation", func(t *testing.T) {
		t.Run("Should return validation error when source is empty", func(t *testing.T) {
			_, err := commandService.CreateCommand(context.Background(), command.CreateCommandParams{
				Source: "",
				Inputs: command.StopMovementInputs{},
			})
			require.Error(t, err)
		})

		t.Run("Should return validation error when inputs are empty", func(t *testing.T) {
			_, err := commandService.CreateCommand(context.Background(), command.CreateCommandParams{
				Source: command.SourceApp,
				Inputs: nil,
			})
			require.Error(t, err)
		})
	})
}

func TestService_CancelCurrentProcessingCommand(t *testing.T) {
	t.Run("Cancel current processing command successfully", func(t *testing.T) {
		runningCommandRepository := commandmocks.NewFakeRunningCommandRepository(t)
		commandService := Service{
			runningCmdRepository: runningCommandRepository,
		}

		cancelableCommand := command.NewCancelableCommand(context.Background(), command.Command{})
		runningCommandRepository.EXPECT().Get(mock.Anything).Return(cancelableCommand, nil)

		err := commandService.CancelCurrentProcessingCommand(context.Background())
		require.NoError(t, err)

		select {
		case <-cancelableCommand.Context().Done():
		case <-time.After(10 * time.Millisecond):
			require.Fail(t, "command should be canceled")
		}
	})

	t.Run("Should return error no command being processed", func(t *testing.T) {
		runningCommandRepository := commandmocks.NewFakeRunningCommandRepository(t)
		commandService := Service{
			runningCmdRepository: runningCommandRepository,
		}

		runningCommandRepository.EXPECT().Get(mock.Anything).Return(command.CancelableCommand{}, command.ErrRunningCommandNotFound)

		err := commandService.CancelCurrentProcessingCommand(context.Background())
		require.Error(t, err)
		require.ErrorIs(t, err, command.ErrNoCommandBeingProcessed)
	})
}

func TestService_ExecuteCreatedCommand(t *testing.T) {
	t.Run("Execute created command successfully", func(t *testing.T) {
		log := logging.NewNoopLogger()
		publisher := eventbusmocks.NewFakePublisher(t)
		runningCommandRepository := commandmocks.NewFakeRunningCommandRepository(t)
		commandRepository := commandmocks.NewFakeRepository(t)
		processingLock := commandmocks.NewFakeProcessingLock(t)
		executorService := commandmocks.NewFakeExecutorService(t)
		commandService := Service{
			log:                  log,
			validator:            validator.New(),
			publisher:            publisher,
			commandRepository:    commandRepository,
			runningCmdRepository: runningCommandRepository,
			processingLock:       processingLock,
			executorService:      executorService,
		}

		runningCommandRepository.EXPECT().Get(mock.Anything).Return(command.CancelableCommand{}, assert.AnError)
		commandRepository.EXPECT().GetCommandByID(mock.Anything, mock.Anything).Return(command.Command{Status: command.StatusQueued}, nil)
		processingLock.EXPECT().WaitUntilUnlocked(mock.Anything).Return(nil)
		executorService.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil)

		doneCh := make(chan struct{})
		commandRepository.EXPECT().GetNextExecutableCommand(mock.Anything).
			Return(command.Command{}, nil).
			Once()
		commandRepository.EXPECT().GetNextExecutableCommand(mock.Anything).
			Return(command.Command{}, command.ErrNoNextExecutableCommand).
			Run(func(_ context.Context) {
				close(doneCh)
			}).
			Once()

		err := commandService.ExecuteCreatedCommand(context.Background(), command.ExecuteCreatedCommandParams{
			CommandID: 1,
		})
		require.NoError(t, err)

		select {
		case <-doneCh:
		case <-time.After(250 * time.Millisecond):
			require.Fail(t, "command should be executed")
		}
	})
}
