package executor

import (
	"context"
	"log/slog"
	"time"

	"github.com/tbe-team/raybot/internal/services/command"
)

func newWaitExecutor(
	log *slog.Logger,
	commandRepository command.Repository,
) *commandExecutor[command.WaitInputs, command.WaitOutputs] {
	handler := waitHandler{
		log: log,
	}

	return newCommandExecutor(
		handler.Handle,
		Hooks[command.WaitOutputs]{},
		log,
		commandRepository,
	)
}

type waitHandler struct {
	log *slog.Logger
}

func (h waitHandler) Handle(ctx context.Context, inputs command.WaitInputs) (command.WaitOutputs, error) {
	select {
	case <-time.After(time.Duration(inputs.DurationMs) * time.Millisecond):
		return command.WaitOutputs{}, nil
	case <-ctx.Done():
		return command.WaitOutputs{}, ctx.Err()
	}
}
