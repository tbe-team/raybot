package command

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/pubsub"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/storage/db"
)

type MoveToLocationExecutor struct {
	commandRepo repository.CommandRepository
	subscriber  message.Subscriber
	picService  service.PICService
	dbProvider  db.Provider
	log         *slog.Logger
}

func NewMoveToLocationExecutor(
	commandRepo repository.CommandRepository,
	subscriber message.Subscriber,
	createSerialCommander service.PICService,
	dbProvider db.Provider,
	log *slog.Logger,
) *MoveToLocationExecutor {
	return &MoveToLocationExecutor{
		commandRepo: commandRepo,
		subscriber:  subscriber,
		picService:  createSerialCommander,
		dbProvider:  dbProvider,
		log:         log,
	}
}

func (e MoveToLocationExecutor) Execute(ctx context.Context, command model.Command) error {
	if command.Type != model.CommandTypeMoveToLocation {
		return fmt.Errorf("command type is not move to location")
	}

	inputs, ok := command.Inputs.(model.CommandMoveToLocationInputs)
	if !ok {
		return fmt.Errorf("command inputs is not model.CommandMoveToLocationInputs")
	}

	execCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start tracking location
	wg := sync.WaitGroup{}
	errChan := make(chan error)
	wg.Add(1)
	go e.trackingLocationLoop(execCtx, &wg, errChan, inputs.Location)

	// Create serial command, robot move forward with speed 100
	params := service.CreateSerialCommandParams{
		Data: model.PICSerialCommandBatteryDriveMotorData{
			Direction: model.MoveDirectionForward,
			Speed:     100,
			Enable:    true,
		},
	}
	if err := e.picService.CreateSerialCommand(ctx, params); err != nil {
		// Create serial command failed, we need to cancel the tracking location loop
		cancel()
		return fmt.Errorf("create serial command: %w", err)
	}

	// Wait for location tracked event
	wg.Wait()

	defer func() {
		// We must stop the robot
		params := service.CreateSerialCommandParams{
			Data: model.PICSerialCommandBatteryDriveMotorData{
				Direction: model.MoveDirectionForward,
				Speed:     0,
				Enable:    true,
			},
		}
		if err := e.picService.CreateSerialCommand(ctx, params); err != nil {
			e.log.Error(
				"robot move to location done but can not stop robot: can not create serial command",
				slog.Any("error", err),
			)
		}
	}()

	// Check if there is an error from the tracking location loop
	select {
	case err := <-errChan:
		return err
	default:
	}

	return nil
}

func (e MoveToLocationExecutor) trackingLocationLoop(
	ctx context.Context,
	wg *sync.WaitGroup,
	errChan chan<- error,
	targetLocation string,
) {
	defer func() {
		wg.Done()
		e.log.Debug("stop tracking location loop")
	}()

	e.log.Debug("start tracking location loop")
	msgChan, err := e.subscriber.Subscribe(ctx, pubsub.TopicRobotLocationUpdated)
	if err != nil {
		errChan <- fmt.Errorf("subscriber subscribe: %w", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		case msg := <-msgChan:
			var event pubsub.RobotLocationUpdatedEvent
			if err := json.Unmarshal(msg.Payload, &event); err != nil {
				errChan <- fmt.Errorf("unmarshal location tracked event: %w", err)
				return
			}

			if event.Location == targetLocation {
				e.log.Debug("current location is the same as the target location")

				// command success, we just return without error
				return
			}

			msg.Ack()
		}
	}
}
