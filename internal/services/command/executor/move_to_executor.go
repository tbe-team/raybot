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

func newMoveToExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	driveMotorService drivemotor.Service,
	commandRepository command.Repository,
) *commandExecutor[command.MoveToInputs] {
	return newCommandExecutor(
		func(ctx context.Context, inputs command.MoveToInputs) error {
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				trackingLocationUntilReached(ctx, inputs.Location, log, subscriber)
			}()

			if err := driveMotorService.MoveForward(ctx, drivemotor.MoveForwardParams{
				Speed: 100,
			}); err != nil {
				return fmt.Errorf("failed to move forward: %w", err)
			}

			wg.Wait()

			if err := driveMotorService.Stop(ctx); err != nil {
				return fmt.Errorf("failed to stop drive motor: %w", err)
			}

			return nil
		},
		Hooks{
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

func trackingLocationUntilReached(
	ctx context.Context,
	location string,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		log.Debug("stop tracking location", slog.String("location", location))
		cancel()
	}()

	doneCh := make(chan struct{})
	log.Debug("start tracking location", slog.String("target_location", location))
	subscriber.Subscribe(ctx, events.LocationUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateLocationEvent)
		if !ok {
			log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if ev.Location == location {
			log.Debug("location reached", slog.String("location", ev.Location))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
