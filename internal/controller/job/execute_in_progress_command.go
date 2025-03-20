package job

import (
	"context"

	"github.com/tbe-team/raybot/internal/service"
)

type ExecuteInProgressCommand struct {
	commandService service.CommandService
}

func NewExecuteInProgressCommand(commandService service.CommandService) *ExecuteInProgressCommand {
	return &ExecuteInProgressCommand{commandService: commandService}
}

func (t ExecuteInProgressCommand) Run(ctx context.Context) error {
	return t.commandService.ExecuteInProgressCommand(ctx)
}
