package event

import (
	"context"

	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/storage/mq"
)

type CommandCreatedEventHandler struct {
	commandService service.CommandService
}

func NewCommandCreatedEventHandler(commandService service.CommandService) *CommandCreatedEventHandler {
	return &CommandCreatedEventHandler{commandService: commandService}
}

func (t CommandCreatedEventHandler) Handle(ctx context.Context, cmdCreatedEvent mq.CommandCreatedEvent) error {
	return t.commandService.ExecuteInProgressCommand(ctx, service.ExecuteInProgressCommandParams{
		CommandID: cmdCreatedEvent.CommandID,
	})
}
