package executor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	configservice "github.com/tbe-team/raybot/internal/services/config"
	"github.com/tbe-team/raybot/internal/services/liftmotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type cargoLowerExecutor struct {
	log              *slog.Logger
	subscriber       eventbus.Subscriber
	configService    configservice.Service
	liftMotorService liftmotor.Service
}

func newCargoLowerExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	configService configservice.Service,
	liftMotorService liftmotor.Service,
) CommandExecutor[command.CargoLowerInputs, command.CargoLowerOutputs] {
	return cargoLowerExecutor{
		log:              log,
		subscriber:       subscriber,
		configService:    configService,
		liftMotorService: liftMotorService,
	}
}

func (e cargoLowerExecutor) Execute(ctx context.Context, inputs command.CargoLowerInputs) (command.CargoLowerOutputs, error) {
	wg := sync.WaitGroup{}

	obstacleCtx, cancelObstacleTracking := context.WithCancel(ctx)
	defer cancelObstacleTracking()

	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingBottomObstacle(obstacleCtx, inputs)
	}()

	readyCh := make(chan struct{}, 1)

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			cancelObstacleTracking()
		}()
		e.trackingLowerPositionUntilReached(ctx, inputs.Position, readyCh)
	}()

	<-readyCh

	if err := e.liftMotorService.SetCargoPosition(ctx, liftmotor.SetCargoPositionParams{
		MotorSpeed: inputs.MotorSpeed,
		Position:   inputs.Position,
	}); err != nil {
		return command.CargoLowerOutputs{}, fmt.Errorf("failed to set cargo position: %w", err)
	}

	// wait for tracking to finish
	wg.Wait()

	if err := e.liftMotorService.Stop(ctx); err != nil {
		return command.CargoLowerOutputs{}, fmt.Errorf("failed to stop lift motor: %w", err)
	}

	return command.CargoLowerOutputs{}, nil
}

func (e cargoLowerExecutor) OnCancel(ctx context.Context) error {
	if err := e.liftMotorService.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop lift motor: %w", err)
	}
	return nil
}

func (e cargoLowerExecutor) trackingLowerPositionUntilReached(ctx context.Context, lowerPosition uint16, readyCh chan<- struct{}) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Info("stop tracking lower position")
		cancel()
	}()

	requiredStableReadCount := e.getRequiredStableReadCount(ctx)
	stableReadCount := 0

	e.log.Info("start tracking lower position",
		slog.Int64("lower_position", int64(lowerPosition)),
		slog.Int("required_stable_read_count", requiredStableReadCount))

	doneCh := make(chan struct{}, 1)
	e.subscriber.Subscribe(ctx, events.DistanceSensorUpdatedTopic, func(msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateDistanceSensorEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if e.isLowerPositionReached(ev.DownDistance, lowerPosition) {
			stableReadCount++
			e.log.Info("lower position reached",
				slog.Int("stable_read_count", stableReadCount),
				slog.Int("required_stable_read_count", requiredStableReadCount),
				slog.Uint64("down_distance", uint64(ev.DownDistance)),
				slog.Int64("lower_position", int64(lowerPosition)))

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

// trackingBottomObstacle tracks the bottom obstacle and stops the motor if it is detected.
// It also starts the motor again if the obstacle is cleared.
// Cancel the context to stop the tracking.
func (e cargoLowerExecutor) trackingBottomObstacle(ctx context.Context, inputs command.CargoLowerInputs) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Info("stop tracking bottom obstacle")
		cancel()
	}()

	obstacleTracking, err := e.getObstacleTracking(ctx)
	if err != nil {
		e.log.Error("failed to get obstacle tracking", slog.Any("error", err))
		return
	}

	bottomDistanceCh := make(chan uint16, 1)

	e.log.Info("start tracking bottom obstacle")
	e.subscriber.Subscribe(ctx, events.CargoBottomDistanceUpdatedTopic, func(msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.CargoBottomDistanceUpdatedEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		select {
		case bottomDistanceCh <- ev.BottomDistance:
		default:
			e.log.Error("dropped message from bottom distance channel",
				slog.Uint64("bottom_distance", uint64(ev.BottomDistance)))
		}
	})

	isMotorRunning := true

	for {
		select {
		case <-ctx.Done():
			return

		case bottomDistance := <-bottomDistanceCh:
			// If the bottom distance is less than the enter distance, we stop the motor
			if bottomDistance <= obstacleTracking.EnterDistance && isMotorRunning {
				e.log.Info("obstacle detected, stopping motor", slog.Uint64("bottom_distance", uint64(bottomDistance)))
				if err := e.liftMotorService.Stop(ctx); err != nil {
					e.log.Error("failed to stop lift motor", slog.Any("error", err))
				}

				isMotorRunning = false
				continue
			}

			// If the bottom distance is greater than the exit distance, we run motor again
			if bottomDistance >= obstacleTracking.ExitDistance && !isMotorRunning {
				e.log.Info("obstacle cleared, running motor again", slog.Uint64("bottom_distance", uint64(bottomDistance)))
				if err := e.liftMotorService.SetCargoPosition(ctx, liftmotor.SetCargoPositionParams{
					MotorSpeed: inputs.MotorSpeed,
					Position:   inputs.Position,
				}); err != nil {
					e.log.Error("failed to set cargo position", slog.Any("error", err))
				}

				isMotorRunning = true
			}
		}
	}
}

func (e cargoLowerExecutor) isLowerPositionReached(current, target uint16) bool {
	acceptableDistance := target - target*10/100 // 10% tolerance
	return current >= acceptableDistance
}

func (e cargoLowerExecutor) getRequiredStableReadCount(ctx context.Context) int {
	commandCfg, err := e.configService.GetCommandConfig(ctx)
	if err != nil {
		e.log.Error("failed to get command config", slog.Any("error", err))
		return 1
	}
	return int(commandCfg.CargoLower.StableReadCount)
}

func (e cargoLowerExecutor) getObstacleTracking(ctx context.Context) (config.ObstacleTracking, error) {
	commandCfg, err := e.configService.GetCommandConfig(ctx)
	if err != nil {
		e.log.Error("failed to get command config", slog.Any("error", err))
		return config.ObstacleTracking{
			EnterDistance: 10,
			ExitDistance:  20,
		}, nil
	}
	return commandCfg.CargoLower.BottomObstacleTracking, nil
}
