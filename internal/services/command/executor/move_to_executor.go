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

const defaultMoveToSpeed = 100

func newMoveToExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	driveMotorService drivemotor.Service,
	commandRepository command.Repository,
) *commandExecutor[command.MoveToInputs, command.MoveToOutputs] {
	handler := moveToHandler{
		log:               log,
		subscriber:        subscriber,
		driveMotorService: driveMotorService,
	}

	return newCommandExecutor(
		handler.Handle,
		Hooks[command.MoveToOutputs]{
			OnCancel: func(ctx context.Context) {
				if err := driveMotorService.Stop(ctx); err != nil {
					log.Error("failed to stop drive motor", slog.Any("error", err))
				}
			},
		},
		log,
		commandRepository,
	)
}

type moveToHandler struct {
	log               *slog.Logger
	subscriber        eventbus.Subscriber
	driveMotorService drivemotor.Service
}

func (h moveToHandler) Handle(ctx context.Context, inputs command.MoveToInputs) (command.MoveToOutputs, error) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		h.trackingLocationUntilReached(ctx, inputs.Location)
	}()

	switch inputs.Direction {
	case command.MoveDirectionForward:
		if err := h.driveMotorService.MoveForward(ctx, drivemotor.MoveForwardParams{
			Speed: defaultMoveToSpeed,
		}); err != nil {
			return command.MoveToOutputs{}, fmt.Errorf("failed to move forward: %w", err)
		}

	case command.MoveDirectionBackward:
		if err := h.driveMotorService.MoveBackward(ctx, drivemotor.MoveBackwardParams{
			Speed: defaultMoveToSpeed,
		}); err != nil {
			return command.MoveToOutputs{}, fmt.Errorf("failed to move backward: %w", err)
		}

	default:
		return command.MoveToOutputs{}, fmt.Errorf("invalid move direction: %s", inputs.Direction)
	}

	wg.Wait()

	if err := h.driveMotorService.Stop(ctx); err != nil {
		return command.MoveToOutputs{}, fmt.Errorf("failed to stop drive motor: %w", err)
	}

	return command.MoveToOutputs{}, nil
}

func (h moveToHandler) trackingLocationUntilReached(ctx context.Context, location string) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		h.log.Debug("stop tracking location", slog.String("location", location))
		cancel()
	}()

	doneCh := make(chan struct{})
	h.log.Debug("start tracking location", slog.String("target_location", location))
	h.subscriber.Subscribe(ctx, events.LocationUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateLocationEvent)
		if !ok {
			h.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if ev.Location == location {
			h.log.Debug("location reached", slog.String("location", ev.Location))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
