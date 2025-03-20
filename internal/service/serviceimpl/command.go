package serviceimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/pkg/validator"
	"github.com/tbe-team/raybot/pkg/xerror"
)

var ErrRobotIsProcessingCommand = xerror.Conflict(nil, "command.alreadyProcessing", "robot is already processing another command")

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

func (s CommandService) CreateCommand(ctx context.Context, params service.CreateCommandParams) (model.Command, error) {
	if err := s.validator.Validate(params); err != nil {
		return model.Command{}, err
	}

	// validate command inputs
	switch params.CommandType {
	case model.CommandTypeMoveToLocation:
		if _, ok := params.Inputs.(model.CommandMoveToLocationInputs); !ok {
			return model.Command{}, xerror.ValidationFailed(nil, "invalid command inputs")
		}
	default:
		return model.Command{}, xerror.ValidationFailed(nil, "invalid command type")
	}

	// check if robot is already processing another command
	if _, err := s.commandRepo.GetCommandByStatusInProgress(ctx, s.dbProvider.DB()); err != nil {
		if !db.IsNoRowsError(err) {
			return model.Command{}, fmt.Errorf("command repository get command by status in progress: %w", err)
		}
	} else {
		return model.Command{}, ErrRobotIsProcessingCommand
	}

	id, err := uuid.NewV7()
	if err != nil {
		return model.Command{}, fmt.Errorf("uuid new v7: %w", err)
	}

	command := model.Command{
		ID:        id.String(),
		Source:    params.Source,
		Type:      params.CommandType,
		Status:    model.CommandStatusInProgress,
		Inputs:    params.Inputs,
		CreatedAt: time.Now(),
	}
	if err := s.commandRepo.CreateCommand(ctx, s.dbProvider.DB(), command); err != nil {
		return model.Command{}, fmt.Errorf("command repository create command: %w", err)
	}

	return command, nil
}
