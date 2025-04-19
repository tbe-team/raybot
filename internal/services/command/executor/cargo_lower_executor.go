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

func newCargoLowerExecutor(
	cfg config.Cargo,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	liftMotorService liftmotor.Service,
	commandRepository command.Repository,
) *commandExecutor[command.CargoLowerInputs, command.CargoLowerOutputs] {
	return newCommandExecutor(
		func(ctx context.Context, _ command.CargoLowerInputs) (command.CargoLowerOutputs, error) {
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				trackingLowerPositionUntilReached(ctx, cfg.LowerPosition, log, subscriber)
			}()

			if err := liftMotorService.SetCargoPosition(ctx, liftmotor.SetCargoPositionParams{
				Position: cfg.LowerPosition,
			}); err != nil {
				return command.CargoLowerOutputs{}, fmt.Errorf("failed to set cargo position: %w", err)
			}

			// wait for distance tracking to finish
			wg.Wait()

			return command.CargoLowerOutputs{}, nil
		},
		Hooks[command.CargoLowerOutputs]{},
		log,
		commandRepository,
	)
}

func trackingLowerPositionUntilReached(
	ctx context.Context,
	lowerPosition uint16,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		log.Debug("stop tracking lower position")
		cancel()
	}()

	doneCh := make(chan struct{})
	log.Debug("start tracking lower position", slog.Int64("lower_position", int64(lowerPosition)))
	subscriber.Subscribe(ctx, events.DistanceSensorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateDistanceSensorEvent)
		if !ok {
			log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		// 10% tolerance
		acceptableDistance := lowerPosition - lowerPosition*10/100
		if ev.DownDistance >= acceptableDistance {
			log.Debug("lower position reached", slog.Int64("lower_position", int64(lowerPosition)))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
