package serviceimpl

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/validator"
	"github.com/tbe-team/raybot/pkg/xerror"
)

var ErrRobotIsProcessingCommand = xerror.Conflict(nil, "command.alreadyProcessing", "robot is already processing another command")

type CreateSerialCommander interface {
	CreateSerialCommand(ctx context.Context, params service.CreateSerialCommandParams) error
}

type CommandService struct {
	commandRepo           repository.CommandRepository
	locationRepo          repository.LocationRepository
	createSerialCommander CreateSerialCommander
	dbProvider            db.Provider
	validator             validator.Validator
	log                   *slog.Logger
}

func NewCommandService(
	commandRepo repository.CommandRepository,
	locationRepo repository.LocationRepository,
	createSerialCommander CreateSerialCommander,
	dbProvider db.Provider,
	validator validator.Validator,
	log *slog.Logger,
) *CommandService {
	return &CommandService{
		commandRepo:           commandRepo,
		locationRepo:          locationRepo,
		createSerialCommander: createSerialCommander,
		dbProvider:            dbProvider,
		validator:             validator,
		log:                   log,
	}
}

func (s CommandService) ListCommands(ctx context.Context, params service.ListCommandsParams) (paging.List[model.Command], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[model.Command]{}, err
	}

	ret, err := s.commandRepo.ListCommands(ctx, s.dbProvider.DB(), params.PagingParams, params.Sorts)
	if err != nil {
		return paging.List[model.Command]{}, fmt.Errorf("command repository list commands: %w", err)
	}

	return ret, nil
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
		if !xerror.IsNotFound(err) {
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

func (s CommandService) ExecuteInProgressCommand(ctx context.Context) error {
	command, err := s.commandRepo.GetCommandByStatusInProgress(ctx, s.dbProvider.DB())
	if err != nil {
		if xerror.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("command repository get command by status in progress: %w", err)
	}

	//nolint:gocritic // it's ok to use switch here we will add more command types in the future
	switch command.Type {
	case model.CommandTypeMoveToLocation:
		inputs, ok := command.Inputs.(model.CommandMoveToLocationInputs)
		if !ok {
			return fmt.Errorf("command inputs is not model.CommandMoveToLocationInputs")
		}
		// We need to start tracking the location
		go func() {
			for {
				loc, err := s.locationRepo.GetCurrentLocation(ctx, s.dbProvider.DB())
				if err != nil {
					s.log.Error("location repository get current location", slog.Any("error", err))
					continue
				}

				if loc.CurrentLocation == inputs.Location {
					s.log.Debug("current location is the same as the target location")
					// Update command status to completed
					completedAt := time.Now()
					params := repository.UpdateCommandParams{
						ID:             command.ID,
						Status:         model.CommandStatusSucceeded,
						SetStatus:      true,
						CompletedAt:    &completedAt,
						SetCompletedAt: true,
					}
					if _, err := s.commandRepo.UpdateCommand(ctx, s.dbProvider.DB(), params); err != nil {
						s.log.Error("command repository update command", slog.Any("error", err))
					}

					break
				}

				time.Sleep(100 * time.Millisecond)
			}
		}()

		if err := s.createSerialCommander.CreateSerialCommand(ctx, service.CreateSerialCommandParams{
			Data: model.PICSerialCommandBatteryDriveMotorData{
				Direction: model.MoveDirectionForward,
				Speed:     100,
				Enable:    true,
			},
		}); err != nil {
			return fmt.Errorf("create serial command: %w", err)
		}
	}

	return nil
}
