package executor

import (
	"context"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type cargoLiftExecutor struct {
	log          *slog.Logger
	liftPosition uint16

	subscriber       eventbus.Subscriber
	liftMotorService liftmotor.Service
}

func newCargoLiftExecutor(
	log *slog.Logger,
	liftPosition uint16,
	subscriber eventbus.Subscriber,
	liftMotorService liftmotor.Service,
) cargoLiftExecutor {
	return cargoLiftExecutor{
		log:              log,
		liftPosition:     liftPosition,
		subscriber:       subscriber,
		liftMotorService: liftMotorService,
	}
}

func (e cargoLiftExecutor) Execute(ctx context.Context, _ command.CargoLiftInputs) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingDistance(ctx)
	}()

	if err := e.liftMotorService.SetCargoPosition(ctx, liftmotor.SetCargoPositionParams{
		Position: e.liftPosition,
	}); err != nil {
		return NewExecutorError(err, "failed to set cargo position")
	}

	// wait for distance tracking to finish
	wg.Wait()

	return nil
}

func (e cargoLiftExecutor) trackingDistance(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Debug("stop tracking distance")
		cancel() // cancel the context to stop the subscriber
	}()

	doneCh := make(chan struct{})
	e.subscriber.Subscribe(ctx, events.DistanceSensorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateDistanceSensorEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		// 10% tolerance
		acceptableDistance := e.liftPosition + e.liftPosition*10/100
		if ev.DownDistance <= acceptableDistance {
			e.log.Debug("cargo lift completed", slog.Int64("lift_position", int64(e.liftPosition)))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
