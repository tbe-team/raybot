package executor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type moveToExecutor struct {
	log               *slog.Logger
	subscriber        eventbus.Subscriber
	driveMotorService drivemotor.Service
}

func newMoveToExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	driveMotorService drivemotor.Service,
) CommandExecutor[command.MoveToInputs, command.MoveToOutputs] {
	return moveToExecutor{
		log:               log,
		subscriber:        subscriber,
		driveMotorService: driveMotorService,
	}
}

func (e moveToExecutor) Execute(ctx context.Context, inputs command.MoveToInputs) (command.MoveToOutputs, error) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingLocationUntilReached(ctx, inputs.Location)
	}()

	switch inputs.Direction {
	case command.MoveDirectionForward:
		if err := e.driveMotorService.MoveForward(ctx, drivemotor.MoveForwardParams{
			Speed: inputs.MotorSpeed,
		}); err != nil {
			return command.MoveToOutputs{}, fmt.Errorf("failed to move forward: %w", err)
		}

	case command.MoveDirectionBackward:
		if err := e.driveMotorService.MoveBackward(ctx, drivemotor.MoveBackwardParams{
			Speed: inputs.MotorSpeed,
		}); err != nil {
			return command.MoveToOutputs{}, fmt.Errorf("failed to move backward: %w", err)
		}

	default:
		return command.MoveToOutputs{}, fmt.Errorf("invalid move direction: %s", inputs.Direction)
	}

	wg.Wait()

	if err := e.driveMotorService.Stop(ctx); err != nil {
		return command.MoveToOutputs{}, fmt.Errorf("failed to stop drive motor: %w", err)
	}

	return command.MoveToOutputs{}, nil
}

func (e moveToExecutor) OnCancel(ctx context.Context) error {
	if err := e.driveMotorService.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop drive motor: %w", err)
	}
	return nil
}

func (e moveToExecutor) trackingLocationUntilReached(ctx context.Context, location string) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Debug("stop tracking location", slog.String("location", location))
		cancel()
	}()

	doneCh := make(chan struct{})
	e.log.Debug("start tracking location", slog.String("target_location", location))
	e.subscriber.Subscribe(ctx, events.LocationUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateLocationEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if ev.Location == location {
			e.log.Debug("location reached", slog.String("location", ev.Location))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
