package executor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/pkg/ptr"
)

type Hooks[O command.Outputs] struct {
	OnSuccess func(ctx context.Context, outputs O)
	OnError   func(ctx context.Context, err error)
	OnCancel  func(ctx context.Context)
}

type commandExecutor[I command.Inputs, O command.Outputs] struct {
	executeFunc func(ctx context.Context, inputs I) (O, error)

	onSuccess func(ctx context.Context, outputs O)
	onError   func(ctx context.Context, err error)
	onCancel  func(ctx context.Context)

	log               *slog.Logger
	commandRepository command.Repository
}

func newCommandExecutor[I command.Inputs, O command.Outputs](
	executeFunc func(ctx context.Context, inputs I) (O, error),
	hooks Hooks[O],
	log *slog.Logger,
	commandRepository command.Repository,
) *commandExecutor[I, O] {
	return &commandExecutor[I, O]{
		executeFunc:       executeFunc,
		onSuccess:         hooks.OnSuccess,
		onError:           hooks.OnError,
		onCancel:          hooks.OnCancel,
		log:               log,
		commandRepository: commandRepository,
	}
}

func (e *commandExecutor[I, O]) Execute(ctx context.Context, cmdID int64, inputs I) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cleanup all resources

	outputs, err := e.executeFunc(ctx, inputs)
	switch {
	case err == nil:
		if e.onSuccess != nil {
			e.onSuccess(ctx, outputs)
		}
		e.updateCommandStatus(cmdID, command.StatusSucceeded, outputs, nil)

	case errors.Is(err, context.Canceled):
		if e.onCancel != nil {
			e.onCancel(ctx)
		}
		e.updateCommandStatus(cmdID, command.StatusCanceled, outputs, nil)

	default:
		if e.onError != nil {
			e.onError(ctx, err)
		}
		e.updateCommandStatus(cmdID, command.StatusFailed, outputs, ptr.New(err.Error()))
	}
}

func (e *commandExecutor[I, O]) updateCommandStatus(
	id int64,
	status command.Status,
	outputs O,
	errMsg *string,
) {
	now := time.Now()
	_, err := e.commandRepository.UpdateCommand(context.TODO(), command.UpdateCommandParams{
		ID:             id,
		Status:         status,
		SetStatus:      true,
		Outputs:        outputs,
		SetOutputs:     true,
		Error:          errMsg,
		SetError:       errMsg != nil,
		CompletedAt:    ptr.New(now),
		SetCompletedAt: true,
		UpdatedAt:      now,
	})
	if err != nil {
		e.log.Error("failed to update command status",
			slog.String("command_id", fmt.Sprintf("%d", id)),
			slog.String("status", string(status)),
			slog.Any("error", err),
		)
	}
}
