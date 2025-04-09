package commandimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/validator"
)

type service struct {
	validator validator.Validator

	commandRepository command.Repository
}

func NewService(
	validator validator.Validator,
	commandRepository command.Repository,
) command.Service {
	return &service{
		validator:         validator,
		commandRepository: commandRepository,
	}
}

func (s service) ListCommands(ctx context.Context, params command.ListCommandsParams) (paging.List[command.Command], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[command.Command]{}, fmt.Errorf("validate params: %w", err)
	}

	return s.commandRepository.ListCommands(ctx, params)
}

func (s service) CreateCommand(ctx context.Context, params command.CreateCommandParams) (command.Command, error) {
	if err := s.validator.Validate(params); err != nil {
		return command.Command{}, fmt.Errorf("validate params: %w", err)
	}

	cmd := command.NewCommand(params.Source, params.Inputs)

	cmd, err := s.commandRepository.CreateCommand(ctx, cmd)
	if err != nil {
		return command.Command{}, fmt.Errorf("create command: %w", err)
	}

	return cmd, nil
}
