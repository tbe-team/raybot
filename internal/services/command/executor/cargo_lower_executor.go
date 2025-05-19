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

type cargoLowerExecutor struct {
	cfg              config.Cargo
	log              *slog.Logger
	subscriber       eventbus.Subscriber
	liftMotorService liftmotor.Service
}

func newCargoLowerExecutor(
	cfg config.Cargo,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	liftMotorService liftmotor.Service,
) CommandExecutor[command.CargoLowerInputs, command.CargoLowerOutputs] {
	return cargoLowerExecutor{
		cfg:              cfg,
		log:              log,
		subscriber:       subscriber,
		liftMotorService: liftMotorService,
	}
}

func (e cargoLowerExecutor) Execute(ctx context.Context, _ command.CargoLowerInputs) (command.CargoLowerOutputs, error) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingLowerPositionUntilReached(ctx, e.cfg.LowerPosition)
	}()

	if err := e.liftMotorService.SetCargoPosition(ctx, liftmotor.SetCargoPositionParams{
		MotorSpeed: e.cfg.LiftMotorSpeed,
		Position:   e.cfg.LowerPosition,
	}); err != nil {
		return command.CargoLowerOutputs{}, fmt.Errorf("failed to set cargo position: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		e.trackingBottomObstacle(ctx)
	}()

	// wait for tracking to finish
	wg.Wait()

	return command.CargoLowerOutputs{}, nil
}

func (e cargoLowerExecutor) OnCancel(ctx context.Context) error {
	if err := e.liftMotorService.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop lift motor: %w", err)
	}
	return nil
}

func (e cargoLowerExecutor) trackingLowerPositionUntilReached(ctx context.Context, lowerPosition uint16) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Debug("stop tracking lower position")
		cancel()
	}()

	doneCh := make(chan struct{})
	e.log.Debug("start tracking lower position", slog.Int64("lower_position", int64(lowerPosition)))
	e.subscriber.Subscribe(ctx, events.DistanceSensorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateDistanceSensorEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		// 10% tolerance
		acceptableDistance := lowerPosition - lowerPosition*10/100
		if ev.DownDistance >= acceptableDistance {
			e.log.Info("lower position reached", slog.Int64("lower_position", int64(lowerPosition)))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}

// trackingBottomObstacle tracks the bottom obstacle and stops the motor if it is detected.
// It also starts the motor again if the obstacle is cleared.
// Caller must ensure that the motor is running before calling this function and cancel the context to stop the tracking.
func (e cargoLowerExecutor) trackingBottomObstacle(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Debug("stop tracking bottom obstacle")
		cancel()
	}()

	bottomDistanceCh := make(chan uint16, 1)

	e.log.Debug("start tracking bottom obstacle")
	e.subscriber.Subscribe(ctx, events.CargoBottomDistanceUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.CargoBottomDistanceUpdatedEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		select {
		case bottomDistanceCh <- ev.BottomDistance:
		default:
			e.log.Error("dropped message from bottom distance channel", slog.Uint64("bottom_distance", uint64(ev.BottomDistance)))
		}
	})

	isMotorRunning := true

	for {
		select {
		case <-ctx.Done():
			return

		case bottomDistance := <-bottomDistanceCh:
			// If the bottom distance is less than the lower threshold, we stop the motor
			if bottomDistance <= e.cfg.BottomDistanceHysteresis.LowerThreshold && isMotorRunning {
				e.log.Info("obstacle detected, stopping motor", slog.Uint64("bottom_distance", uint64(bottomDistance)))
				if err := e.liftMotorService.Stop(ctx); err != nil {
					e.log.Error("failed to stop lift motor", slog.Any("error", err))
				}

				isMotorRunning = false
				continue
			}

			// If the bottom distance is greater than the upper threshold, we run motor again
			if bottomDistance >= e.cfg.BottomDistanceHysteresis.UpperThreshold && !isMotorRunning {
				e.log.Info("obstacle cleared, running motor again", slog.Uint64("bottom_distance", uint64(bottomDistance)))
				if err := e.liftMotorService.SetCargoPosition(ctx, liftmotor.SetCargoPositionParams{
					MotorSpeed: e.cfg.LiftMotorSpeed,
					Position:   e.cfg.LowerPosition,
				}); err != nil {
					e.log.Error("failed to set cargo position", slog.Any("error", err))
				}

				isMotorRunning = true
			}
		}
	}
}
