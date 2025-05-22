package executor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type cargoLiftExecutor struct {
	log              *slog.Logger
	subscriber       eventbus.Subscriber
	liftMotorService liftmotor.Service
}

func newCargoLiftExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	liftMotorService liftmotor.Service,
) CommandExecutor[command.CargoLiftInputs, command.CargoLiftOutputs] {
	return cargoLiftExecutor{
		log:              log,
		subscriber:       subscriber,
		liftMotorService: liftMotorService,
	}
}

func (e cargoLiftExecutor) Execute(ctx context.Context, inputs command.CargoLiftInputs) (command.CargoLiftOutputs, error) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingLiftPositionUntilReached(ctx, inputs.Position)
	}()

	if err := e.liftMotorService.SetCargoPosition(ctx, liftmotor.SetCargoPositionParams{
		MotorSpeed: inputs.MotorSpeed,
		Position:   inputs.Position,
	}); err != nil {
		return command.CargoLiftOutputs{}, fmt.Errorf("failed to set cargo position: %w", err)
	}

	// wait for tracking to finish
	wg.Wait()

	return command.CargoLiftOutputs{}, nil
}

func (e cargoLiftExecutor) OnCancel(ctx context.Context) error {
	if err := e.liftMotorService.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop lift motor: %w", err)
	}
	return nil
}

func (e cargoLiftExecutor) trackingLiftPositionUntilReached(ctx context.Context, liftPosition uint16) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Debug("stop tracking lift position")
		cancel()
	}()

	doneCh := make(chan struct{})
	e.log.Debug("start tracking lift position", slog.Int64("lift_position", int64(liftPosition)))
	e.subscriber.Subscribe(ctx, events.DistanceSensorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateDistanceSensorEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		// 10% tolerance
		acceptableDistance := liftPosition + liftPosition*10/100
		if ev.DownDistance <= acceptableDistance {
			e.log.Debug("lift position reached", slog.Int64("lift_position", int64(liftPosition)))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
