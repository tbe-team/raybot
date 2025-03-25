package event

import (
	"context"

	"github.com/tbe-team/raybot/internal/pubsub"
	"github.com/tbe-team/raybot/internal/service"
)

type CommandCreatedEventHandler struct {
	commandService service.CommandService
}

func NewCommandCreatedEventHandler(commandService service.CommandService) *CommandCreatedEventHandler {
	return &CommandCreatedEventHandler{commandService: commandService}
}

func (t CommandCreatedEventHandler) Handle(ctx context.Context, cmdCreatedEvent pubsub.CommandCreatedEvent) error {
	return t.commandService.ExecuteCommand(ctx, service.ExecuteCommandParams{
		CommandID: cmdCreatedEvent.CommandID,
	})
}
