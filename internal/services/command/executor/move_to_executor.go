package executor

import (
	"context"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
)

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
	if err := e.driveMotorService.UpdateDriveMotorState(ctx, drivemotor.UpdateDriveMotorStateParams{
		Direction: drivemotor.DirectionForward,
		Speed:     100,
		Enabled:   true,
	}); err != nil {
		return NewExecutorError(err, "failed to update drive motor state, (start driving)")
	}

	// wait for location tracking to finish
	wg.Wait()

	// stop driving
	if err := e.driveMotorService.UpdateDriveMotorState(ctx, drivemotor.UpdateDriveMotorStateParams{
		Direction: drivemotor.DirectionForward,
		Speed:     0,
		Enabled:   true,
	}); err != nil {
		return NewExecutorError(err, "failed to update drive motor state, (stop driving)")
	}

	return nil
}

func (e moveToExecutor) trackingLocation(ctx context.Context, location string) {
	locationTrackingKey := "location:tracking"
	doneCh := make(chan struct{})
	defer events.UpdateLocationSignal.RemoveListener(locationTrackingKey)

	var once sync.Once

	fn := func(_ context.Context, ev events.UpdateLocationEvent) {
		if ev.CurrentLocation == location {
			e.log.Info("location reached", slog.String("location", ev.CurrentLocation))
			once.Do(func() {
				close(doneCh)
			})
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
