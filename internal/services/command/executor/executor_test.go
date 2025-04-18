package executor

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tbe-team/raybot/internal/services/command"
	commandmock "github.com/tbe-team/raybot/internal/services/command/mocks"
	"github.com/tbe-team/raybot/pkg/log"
)

func TestCommandExecutor(t *testing.T) {
	cmdID := int64(1)
	t.Run("Execute successfully", func(t *testing.T) {
		successCalled := false
		cancelCalled := false
		errorCalled := false

		executor, fakeCommandRepo := newTestExecutor(
			t,
			func(_ context.Context, _ command.StopMovementInputs) (command.StopMovementOutputs, error) {
				return command.StopMovementOutputs{}, nil
			},
			Hooks[command.StopMovementOutputs]{
				OnSuccess: func(_ context.Context, _ command.StopMovementOutputs) {
					successCalled = true
				},
				OnCancel: func(_ context.Context) {
					cancelCalled = true
				},
				OnError: func(_ context.Context, _ error) {
					errorCalled = true
				},
			},
		)
		fakeCommandRepo.EXPECT().UpdateCommand(
			mock.Anything,
			mock.MatchedBy(func(p command.UpdateCommandParams) bool {
				return p.ID == cmdID &&
					p.Status == command.StatusSucceeded &&
					p.SetStatus == true &&
					p.Error == nil &&
					p.SetError == false &&
					p.CompletedAt != nil &&
					p.SetCompletedAt == true &&
					!p.UpdatedAt.IsZero()
			}),
		).Return(command.Command{}, nil)
		executor.Execute(context.Background(), cmdID, command.StopMovementInputs{})

		assert.True(t, successCalled)
		assert.False(t, cancelCalled)
		assert.False(t, errorCalled)
		fakeCommandRepo.AssertExpectations(t)
	})

	t.Run("Execute canceled when running", func(t *testing.T) {
		successCalled := false
		cancelCalled := false
		errorCalled := false

		executor, fakeCommandRepo := newTestExecutor(
			t,
			func(ctx context.Context, _ command.StopMovementInputs) (command.StopMovementOutputs, error) {
				<-ctx.Done()
				return command.StopMovementOutputs{}, ctx.Err()
			},
			Hooks[command.StopMovementOutputs]{
				OnSuccess: func(_ context.Context, _ command.StopMovementOutputs) {
					successCalled = true
				},
				OnCancel: func(_ context.Context) {
					cancelCalled = true
				},
				OnError: func(_ context.Context, _ error) {
					errorCalled = true
				},
			},
		)
		fakeCommandRepo.EXPECT().UpdateCommand(
			mock.Anything,
			mock.MatchedBy(func(p command.UpdateCommandParams) bool {
				return p.ID == cmdID &&
					p.Status == command.StatusCanceled &&
					p.SetStatus == true &&
					p.Error == nil &&
					p.SetError == false &&
					p.CompletedAt != nil &&
					p.SetCompletedAt == true &&
					!p.UpdatedAt.IsZero()
			}),
		).Return(command.Command{}, nil)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		executor.Execute(ctx, cmdID, command.StopMovementInputs{})

		assert.False(t, successCalled)
		assert.True(t, cancelCalled)
		assert.False(t, errorCalled)
		fakeCommandRepo.AssertExpectations(t)
	})

	t.Run("Execute failed", func(t *testing.T) {
		successCalled := false
		cancelCalled := false
		errorCalled := false

		failedError := errors.New("failed error")
		executor, fakeCommandRepo := newTestExecutor(
			t,
			func(_ context.Context, _ command.StopMovementInputs) (command.StopMovementOutputs, error) {
				return command.StopMovementOutputs{}, failedError
			},
			Hooks[command.StopMovementOutputs]{
				OnSuccess: func(_ context.Context, _ command.StopMovementOutputs) {
					successCalled = true
				},
				OnCancel: func(_ context.Context) {
					cancelCalled = true
				},
				OnError: func(_ context.Context, _ error) {
					errorCalled = true
				},
			},
		)
		fakeCommandRepo.EXPECT().UpdateCommand(
			mock.Anything,
			mock.MatchedBy(func(p command.UpdateCommandParams) bool {
				return p.ID == cmdID &&
					p.Status == command.StatusFailed &&
					p.SetStatus == true &&
					p.Error != nil &&
					*p.Error == failedError.Error() &&
					p.SetError == true &&
					p.CompletedAt != nil &&
					p.SetCompletedAt == true &&
					!p.UpdatedAt.IsZero()
			}),
		).Return(command.Command{}, nil)
		executor.Execute(context.Background(), cmdID, command.StopMovementInputs{})

		assert.False(t, successCalled)
		assert.False(t, cancelCalled)
		assert.True(t, errorCalled)
		fakeCommandRepo.AssertExpectations(t)
	})
}

func newTestExecutor[I command.Inputs, O command.Outputs](
	t *testing.T,
	executeFunc func(ctx context.Context, inputs I) (O, error),
	hooks Hooks[O],
) (*commandExecutor[I, O], *commandmock.FakeRepository) {
	log := log.NewNoopLogger()
	commandRepo := commandmock.NewFakeRepository(t)
	exec := newCommandExecutor(executeFunc, hooks, log, commandRepo)
	return exec, commandRepo
}
