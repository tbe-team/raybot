package executor

import (
	"context"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
)

const locationTrackingKey = "location:tracking"

type moveToExecutor struct {
	log               *slog.Logger
	driveMotorService drivemotor.Service
}

func newMoveToExecutor(
	log *slog.Logger,
	driveMotorService drivemotor.Service,
) moveToExecutor {
	return moveToExecutor{
		log:               log,
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
	e.log.Debug("start tracking location", slog.String("location", location))

	doneCh := make(chan struct{})

	defer func() {
		e.log.Debug("stop tracking location", slog.String("location", location))
		events.UpdateLocationSignal.RemoveListener(locationTrackingKey)
	}()

	// var once sync.Once

	fn := func(_ context.Context, ev events.UpdateLocationEvent) {
		if ev.CurrentLocation == location {
			e.log.Info("location reached", slog.String("location", ev.CurrentLocation))
			close(doneCh)

		}
	}

	events.UpdateLocationSignal.AddListener(fn, locationTrackingKey)

	select {
	case <-doneCh:
		return
	case <-ctx.Done():
		return
	}
}
