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

type cargoLowerExecutor struct {
	log           *slog.Logger
	lowerPosition uint16

	subscriber       eventbus.Subscriber
	liftMotorService liftmotor.Service
}

func newCargoLowerExecutor(
	log *slog.Logger,
	lowerPosition uint16,
	subscriber eventbus.Subscriber,
	liftMotorService liftmotor.Service,
) cargoLowerExecutor {
	return cargoLowerExecutor{
		log:              log,
		lowerPosition:    lowerPosition,
		subscriber:       subscriber,
		liftMotorService: liftMotorService,
	}
}

func (e cargoLowerExecutor) Execute(ctx context.Context, _ command.CargoLowerInputs) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingDistance(ctx)
	}()

	if err := e.liftMotorService.SetCargoPosition(ctx, liftmotor.SetCargoPositionParams{
		Position: e.lowerPosition,
	}); err != nil {
		return NewExecutorError(err, "failed to set cargo position")
	}

	// wait for distance tracking to finish
	wg.Wait()

	return nil
}

func (e cargoLowerExecutor) trackingDistance(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	doneCh := make(chan struct{})
	e.subscriber.Subscribe(ctx, events.DistanceSensorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateDistanceSensorEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		// 10% tolerance
		acceptableDistance := e.lowerPosition - e.lowerPosition*10/100
		if ev.DownDistance >= acceptableDistance {
			e.log.Debug("cargo lower completed", slog.Int64("lower_position", int64(e.lowerPosition)))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
