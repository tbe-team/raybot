package executor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/config"
	"github.com/tbe-team/raybot/internal/services/distancesensor"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type cargoLiftExecutor struct {
	log                   *slog.Logger
	subscriber            eventbus.Subscriber
	configService         config.Service
	liftMotorService      liftmotor.Service
	distanceSensorService distancesensor.Service
}

func newCargoLiftExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	configService config.Service,
	liftMotorService liftmotor.Service,
	distanceSensorService distancesensor.Service,
) CommandExecutor[command.CargoLiftInputs, command.CargoLiftOutputs] {
	return cargoLiftExecutor{
		log:                   log,
		subscriber:            subscriber,
		configService:         configService,
		liftMotorService:      liftMotorService,
		distanceSensorService: distanceSensorService,
	}
}

func (e cargoLiftExecutor) Execute(ctx context.Context, inputs command.CargoLiftInputs) (command.CargoLiftOutputs, error) {
	distanceSensorState, err := e.distanceSensorService.GetDistanceSensorState(ctx)
	if err != nil {
		return command.CargoLiftOutputs{}, fmt.Errorf("failed to get distance sensor state: %w", err)
	}

	if e.isLiftPositionReached(distanceSensorState.DownDistance, inputs.Position) {
		return command.CargoLiftOutputs{}, nil
	}

	wg := sync.WaitGroup{}
	readyCh := make(chan struct{}, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingLiftPositionUntilReached(ctx, inputs.Position, readyCh)
	}()

	<-readyCh

	if err := e.liftMotorService.SetCargoPosition(ctx, liftmotor.SetCargoPositionParams{
		MotorSpeed: inputs.MotorSpeed,
		Position:   inputs.Position,
	}); err != nil {
		return command.CargoLiftOutputs{}, fmt.Errorf("failed to set cargo position: %w", err)
	}

	// wait for tracking to finish
	wg.Wait()

	if err := e.liftMotorService.Stop(ctx); err != nil {
		return command.CargoLiftOutputs{}, fmt.Errorf("failed to stop lift motor: %w", err)
	}

	return command.CargoLiftOutputs{}, nil
}

func (e cargoLiftExecutor) OnCancel(ctx context.Context) error {
	if err := e.liftMotorService.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop lift motor: %w", err)
	}
	return nil
}

func (e cargoLiftExecutor) trackingLiftPositionUntilReached(ctx context.Context, liftPosition uint16, readyCh chan<- struct{}) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Info("stop tracking lift position")
		cancel()
	}()

	requiredStableReadCount := e.getRequiredStableReadCount(ctx)
	stableReadCount := 0

	e.log.Info("start tracking lift position",
		slog.Int64("target_position", int64(liftPosition)),
		slog.Int("required_stable_read_count", requiredStableReadCount))

	doneCh := make(chan struct{}, 1)
	e.subscriber.Subscribe(ctx, events.DistanceSensorUpdatedTopic, func(msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateDistanceSensorEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if e.isLiftPositionReached(ev.DownDistance, liftPosition) {
			stableReadCount++
			e.log.Info("lift position reached",
				slog.Int("stable_read_count", stableReadCount),
				slog.Int("required_stable_read_count", requiredStableReadCount),
				slog.Int64("down_distance", int64(ev.DownDistance)),
				slog.Int64("target_position", int64(liftPosition)))

			if stableReadCount >= requiredStableReadCount {
				select {
				case doneCh <- struct{}{}:
				default:
				}
			}
			return
		}

		e.log.Warn("reset stable read count",
			slog.Int64("down_distance", int64(ev.DownDistance)))
		stableReadCount = 0
	})

	readyCh <- struct{}{}

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}

func (cargoLiftExecutor) isLiftPositionReached(current, target uint16) bool {
	acceptableDistance := target + target*10/100 // 10% tolerance
	return current <= acceptableDistance
}

func (e cargoLiftExecutor) getRequiredStableReadCount(ctx context.Context) int {
	commandCfg, err := e.configService.GetCommandConfig(ctx)
	if err != nil {
		e.log.Error("failed to get command config", slog.Any("error", err))
		return 1
	}
	return int(commandCfg.CargoLift.StableReadCount)
}
