package executor

import (
	"context"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type moveToExecutor struct {
	log *slog.Logger

	subscriber        eventbus.Subscriber
	driveMotorService drivemotor.Service
}

func newMoveToExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	driveMotorService drivemotor.Service,
) moveToExecutor {
	return moveToExecutor{
		log:               log,
		subscriber:        subscriber,
		driveMotorService: driveMotorService,
	}
}

func (e moveToExecutor) Execute(ctx context.Context, inputs command.MoveToInputs) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingLocation(ctx, inputs.Location)
	}()

	// start driving
	if err := e.driveMotorService.MoveForward(ctx, drivemotor.MoveForwardParams{
		Speed: 100,
	}); err != nil {
		return NewExecutorError(err, "failed to move forward")
	}

	// wait for location tracking to finish
	wg.Wait()

	// stop driving
	if err := e.driveMotorService.Stop(ctx); err != nil {
		return NewExecutorError(err, "failed to stop driving")
	}

	return nil
}

func (e moveToExecutor) trackingLocation(ctx context.Context, location string) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Debug("stop tracking location", slog.String("location", location))
		cancel() // cancel the context to stop the subscriber
	}()

	reachCh := make(chan struct{})
	e.log.Debug("start tracking location", slog.String("target_location", location))
	e.subscriber.Subscribe(ctx, events.LocationUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateLocationEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if ev.Location == location {
			e.log.Debug("location reached", slog.String("location", ev.Location))
			close(reachCh)
		}
	})

	select {
	case <-reachCh:
	case <-ctx.Done():
	}
}
