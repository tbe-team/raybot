package executor

import (
	"context"
	"time"

	"github.com/tbe-team/raybot/internal/services/command"
)

type waitExecutor struct{}

func newWaitExecutor() CommandExecutor[command.WaitInputs, command.WaitOutputs] {
	return waitExecutor{}
}

func (e waitExecutor) Execute(ctx context.Context, inputs command.WaitInputs) (command.WaitOutputs, error) {
	select {
	case <-time.After(time.Duration(inputs.DurationMs) * time.Millisecond):
		return command.WaitOutputs{}, nil
	case <-ctx.Done():
		return command.WaitOutputs{}, ctx.Err()
	}
}

func (e waitExecutor) OnCancel(_ context.Context) error {
	return nil
}
