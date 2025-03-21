package command

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/pubsub"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/validator"
	"github.com/tbe-team/raybot/pkg/xerror"
)

var ErrRobotIsProcessingCommand = xerror.Conflict(nil, "command.alreadyProcessing", "robot is already processing another command")

type CreateSerialServicer interface {
	CreateSerialCommand(ctx context.Context, params service.CreateSerialCommandParams) error
}

type Service struct {
	commandRepo             repository.CommandRepository
	createSerialServicer    CreateSerialServicer
	dbProvider              db.Provider
	publisher               message.Publisher
	subscriber              message.Subscriber
	validator               validator.Validator
	commandExecutorRegistry *registry
	log                     *slog.Logger
}

func NewService(
	commandRepo repository.CommandRepository,
	createSerialCommander CreateSerialServicer,
	dbProvider db.Provider,
	publisher message.Publisher,
	subscriber message.Subscriber,
	validator validator.Validator,
	log *slog.Logger,
) *Service {
	service := &Service{
		commandRepo:             commandRepo,
		createSerialServicer:    createSerialCommander,
		dbProvider:              dbProvider,
		publisher:               publisher,
		subscriber:              subscriber,
		validator:               validator,
		commandExecutorRegistry: newRegistry(),
		log:                     log,
	}

	service.registerCommandExecutors()

	return service
}

func (s Service) ListCommands(ctx context.Context, params service.ListCommandsParams) (paging.List[model.Command], error) {
	if err := s.validator.Validate(params); err != nil {
		return paging.List[model.Command]{}, err
	}

	ret, err := s.commandRepo.ListCommands(ctx, s.dbProvider.DB(), params.PagingParams, params.Sorts)
	if err != nil {
		return paging.List[model.Command]{}, fmt.Errorf("command repository list commands: %w", err)
	}

	return ret, nil
}

func (s Service) GetCurrentProcessingCommand(ctx context.Context) (model.Command, error) {
	command, err := s.commandRepo.GetCommandByStatusInProgress(ctx, s.dbProvider.DB())
	if err != nil {
		return model.Command{}, fmt.Errorf("command repository get current processing command: %w", err)
	}

	return command, nil
}

func (s Service) CreateCommand(ctx context.Context, params service.CreateCommandParams) (model.Command, error) {
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

	if err := s.publishCommandCreatedEvent(command); err != nil {
		return model.Command{}, fmt.Errorf("publish command created event: %w", err)
	}

	return command, nil
}

func (s Service) publishCommandCreatedEvent(command model.Command) error {
	ev := pubsub.CommandCreatedEvent{
		CommandID: command.ID,
	}
	payload, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("json marshal command created event: %w", err)
	}

	msg := message.NewMessage(shortuuid.New(), payload)
	if err := s.publisher.Publish(pubsub.TopicCommandCreated, msg); err != nil {
		return fmt.Errorf("publisher publish command created event: %w", err)
	}

	return nil
}

func (s Service) ExecuteCommand(ctx context.Context, params service.ExecuteCommandParams) error {
	if err := s.validator.Validate(params); err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	command, err := s.commandRepo.GetCommandByID(ctx, s.dbProvider.DB(), params.CommandID)
	if err != nil {
		if xerror.IsNotFound(err) {
			s.log.Error("command not found", slog.Any("command_id", params.CommandID))
			return nil
		}
		return fmt.Errorf("command repository get command by id: %w", err)
	}

	// if command is not in progress, we don't need to execute it
	if command.Status != model.CommandStatusInProgress {
		return nil
	}

	// get executor for command type
	executor, err := s.commandExecutorRegistry.GetExecutor(command.Type)
	if err != nil {
		s.log.Error("no executor found for command type",
			slog.String("command_type", command.Type.String()))
		errorMessage := fmt.Sprintf("no executor found for command type: %s", command.Type)
		params := repository.UpdateCommandParams{
			ID:        command.ID,
			Status:    model.CommandStatusFailed,
			SetStatus: true,
			Error:     &errorMessage,
			SetError:  true,
		}
		if _, err := s.commandRepo.UpdateCommand(ctx, s.dbProvider.DB(), params); err != nil {
			return fmt.Errorf("command repository update command: %w", err)
		}

		return nil
	}

	return executor.Execute(ctx, command)
}

func (s Service) registerCommandExecutors() {
	s.commandExecutorRegistry.Register(
		model.CommandTypeMoveToLocation,
		NewMoveToLocationExecutor(
			s.commandRepo,
			s.subscriber,
			s.createSerialServicer,
			s.dbProvider,
			s.log,
		),
	)
}
