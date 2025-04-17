package executor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

func newCargoLiftExecutor(
	cfg config.Cargo,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	liftMotorService liftmotor.Service,
	commandRepository command.Repository,
) *commandExecutor[command.CargoLiftInputs] {
	return newCommandExecutor(
		func(ctx context.Context, _ command.CargoLiftInputs) error {
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				trackingLiftPositionUntilReached(ctx, cfg.LiftPosition, log, subscriber)
			}()

			if err := liftMotorService.SetCargoPosition(ctx, liftmotor.SetCargoPositionParams{
				Position: cfg.LiftPosition,
			}); err != nil {
				return fmt.Errorf("failed to set cargo position: %w", err)
			}

			// wait for distance tracking to finish
			wg.Wait()

			return nil
		},
		Hooks{},
		log,
		commandRepository,
	)
}

func trackingLiftPositionUntilReached(
	ctx context.Context,
	liftPosition uint16,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		log.Debug("stop tracking lift position")
		cancel()
	}()

	doneCh := make(chan struct{})
	log.Debug("start tracking lift position", slog.Int64("lift_position", int64(liftPosition)))
	subscriber.Subscribe(ctx, events.DistanceSensorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateDistanceSensorEvent)
		if !ok {
			log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		// 10% tolerance
		acceptableDistance := liftPosition + liftPosition*10/100
		if ev.DownDistance <= acceptableDistance {
			log.Debug("lift position reached", slog.Int64("lift_position", int64(liftPosition)))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
