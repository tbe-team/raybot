package serviceimpl

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/pkg/validator"
)

type CommandService struct {
	commandRepo repository.CommandRepository
	dbProvider  db.Provider
	validator   validator.Validator
}

func NewCommandService(commandRepo repository.CommandRepository, dbProvider db.Provider, validator validator.Validator) *CommandService {
	return &CommandService{
		commandRepo: commandRepo,
		dbProvider:  dbProvider,
		validator:   validator,
	}
}

func (s CommandService) ListCommands(ctx context.Context, params service.ListCommandsParams) ([]model.Command, error) {
	if err := s.validator.Validate(params); err != nil {
		return nil, err
	}

	commands, err := s.commandRepo.ListCommands(ctx, s.dbProvider.DB(), params.PagingParams, params.Sorts)
	if err != nil {
		return nil, fmt.Errorf("command repository list commands: %w", err)
	}

	return commands, nil
}

func (s CommandService) GetCurrentProcessingCommand(ctx context.Context) (model.Command, error) {
	command, err := s.commandRepo.GetCommandByStatusInProgress(ctx, s.dbProvider.DB())
	if err != nil {
		return model.Command{}, fmt.Errorf("command repository get current processing command: %w", err)
	}

	return command, nil
}
