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

type stopMovementExecutor struct {
	log *slog.Logger

	subscriber        eventbus.Subscriber
	driveMotorService drivemotor.Service
}

func newStopMovementExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	driveMotorService drivemotor.Service,
) stopMovementExecutor {
	return stopMovementExecutor{
		log:               log,
		subscriber:        subscriber,
		driveMotorService: driveMotorService,
	}
}

func (e stopMovementExecutor) Execute(ctx context.Context, _ command.StopMovementInputs) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingDriveMotorStop(ctx)
	}()

	if err := e.driveMotorService.Stop(ctx); err != nil {
		return NewExecutorError(err, "failed to stop")
	}

	wg.Wait()

	return nil
}

func (e stopMovementExecutor) trackingDriveMotorStop(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Debug("stop tracking drive motor stop")
		cancel() // cancel the context to stop the subscriber
	}()

	doneCh := make(chan struct{})
	e.log.Debug("start tracking drive motor stop")
	e.subscriber.Subscribe(ctx, events.DriveMotorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.DriveMotorStateUpdatedEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if !ev.IsRunning {
			e.log.Debug("drive motor stopped")
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
