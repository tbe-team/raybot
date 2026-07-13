package executor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type moveToExecutor struct {
	log               *slog.Logger
	subscriber        eventbus.Subscriber
	driveMotorService drivemotor.Service
}

func newMoveToExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	driveMotorService drivemotor.Service,
) CommandExecutor[command.MoveToInputs, command.MoveToOutputs] {
	return moveToExecutor{
		log:               log,
		subscriber:        subscriber,
		driveMotorService: driveMotorService,
	}
}

func (e moveToExecutor) Execute(ctx context.Context, inputs command.MoveToInputs) (command.MoveToOutputs, error) {
	wg := sync.WaitGroup{}

	runCtx, cancelRunCtx := context.WithCancel(ctx)
	defer cancelRunCtx()

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			cancelRunCtx() // if reached location, we don't need to run anymore
		}()
		e.trackingLocationUntilReached(ctx, inputs.Location)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		e.run(runCtx, inputs)
	}()

	wg.Wait()

	if err := e.driveMotorService.Stop(ctx); err != nil {
		return command.MoveToOutputs{}, fmt.Errorf("failed to stop drive motor: %w", err)
	}

	return command.MoveToOutputs{}, nil
}

func (e moveToExecutor) OnCancel(ctx context.Context) error {
	if err := e.driveMotorService.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop drive motor: %w", err)
	}
	return nil
}

func (e moveToExecutor) trackingLocationUntilReached(ctx context.Context, location string) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Info("stop tracking location", slog.String("location", location))
		cancel()
	}()

	doneCh := make(chan struct{})
	e.log.Info("start tracking location", slog.String("target_location", location))
	e.subscriber.Subscribe(ctx, events.LocationUpdatedTopic, func(msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.LocationUpdatedEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if ev.Location == location {
			e.log.Info("location reached", slog.String("location", ev.Location))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}

func (e moveToExecutor) run(ctx context.Context, inputs command.MoveToInputs) {
	e.log.Info("start running with safety monitor", slog.Any("inputs", inputs))

	evChan := make(chan events.UpdateDistanceSensorEvent, 1)
	e.subscriber.Subscribe(ctx, events.DistanceSensorUpdatedTopic, func(msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateDistanceSensorEvent)
		if !ok {
			e.log.Error("invalid distance sensor event", slog.Any("event", msg.Payload))
			return
		}

		select {
		case evChan <- ev:
		default:
			e.log.Warn("distance sensor event channel is full, dropping event", slog.Any("event", ev))
		}
	})

	runDriveMotorFunc := func() error {
		// TODO: we can optimize this func
		return e.runDriveMotor(ctx, inputs.Direction, inputs.MotorSpeed)
	}

	stopDriveMotorFunc := func() error {
		return e.driveMotorService.Stop(ctx)
	}

	h := moveToSafetyHandler{
		log:                  e.log,
		stopDriveMotorFunc:   stopDriveMotorFunc,
		resumeDriveMotorFunc: runDriveMotorFunc,
		Direction:            inputs.Direction,
		ObstacleTracking:     inputs.ObstacleTracking,
		CargoLostDistance:    inputs.CargoLostDistance,
	}

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-evChan:
			h.Handle(ev)
		}
	}
}

func (e moveToExecutor) runDriveMotor(ctx context.Context, direction command.MoveDirection, speed uint8) error {
	switch direction {
	case command.MoveDirectionForward:
		if err := e.driveMotorService.MoveForward(ctx, drivemotor.MoveForwardParams{
			Speed: speed,
		}); err != nil {
			return fmt.Errorf("failed to move forward: %w", err)
		}

	case command.MoveDirectionBackward:
		if err := e.driveMotorService.MoveBackward(ctx, drivemotor.MoveBackwardParams{
			Speed: speed,
		}); err != nil {
			return fmt.Errorf("failed to move backward: %w", err)
		}

	default:
		return fmt.Errorf("invalid move direction: %s", direction)
	}

	return nil
}

type moveToSafetyHandler struct {
	log                  *slog.Logger
	stopDriveMotorFunc   func() error
	resumeDriveMotorFunc func() error

	Direction         command.MoveDirection
	ObstacleTracking  command.ObstacleTracking
	CargoLostDistance uint16

	isDriveMotorRunning bool
	cargoLost           bool
	obstacleDetected    bool
}

func (m *moveToSafetyHandler) Handle(ev events.UpdateDistanceSensorEvent) {
	currentCargoLost := m.isCargoLost(ev.DownDistance)

	// cargo lost
	if currentCargoLost && !m.cargoLost {
		m.cargoLost = true
		m.log.Warn("cargo lost detected", slog.Uint64("down_distance", uint64(ev.DownDistance)))
		m.stopMovement()
		// TODO: tracking lift cargo here, should block until cargo is lifted
		return
	}

	// cargo back in position
	if !currentCargoLost && m.cargoLost {
		m.cargoLost = false
		m.log.Info("cargo back in position", slog.Uint64("down_distance", uint64(ev.DownDistance)))
	}

	// If cargo is lost, we don't need to handle obstacle detection
	// cargo lost has higher priority than obstacle detection
	if currentCargoLost {
		return
	}

	currentObstacleDetected := m.isObstacleDetected(ev)

	// obstacle detected
	if currentObstacleDetected && !m.obstacleDetected {
		m.obstacleDetected = true
		m.log.Warn("obstacle detected",
			slog.String("direction", string(m.Direction)),
			slog.Uint64("front_distance", uint64(ev.FrontDistance)),
			slog.Uint64("back_distance", uint64(ev.BackDistance)))
		m.stopMovement()
		return
	}

	// obstacle cleared
	if !currentObstacleDetected && m.obstacleDetected {
		m.obstacleDetected = false
		m.log.Info("obstacle cleared",
			slog.String("direction", string(m.Direction)),
			slog.Uint64("front_distance", uint64(ev.FrontDistance)),
			slog.Uint64("back_distance", uint64(ev.BackDistance)))
		m.resumeMovement()
		return
	}
}

//   - There are two thresholds: EnterDistance (start detecting obstacle),
//     ExitDistance (stop detecting obstacle).
//
//   - If `obstacleDetected` is already true (we are currently tracking an obstacle):
//
//   - Only stop tracking when actual distance > ExitDistance.
//
//   - Otherwise, still considered in obstacle state.
//
//   - If `obstacleDetected` is false:
//
//   - Start tracking when actual distance < EnterDistance.
//
// Diagram:
// `
//
//	   ↓ decreasing distance                  ↑ increasing distance
//	┌────────────────────────────┐      ┌───────────────────────────┐
//	│ actualDistance < EnterDist ├─────▶│ obstacleDetected = true   │
//	└────────────────────────────┘      └───────────────────────────┘
//	                                             │
//	       ┌─────────────────────────────────────┘
//	       │
//	       ▼
//	┌────────────────────────────┐
//	│ actualDistance > ExitDist  ├─────▶ obstacleDetected = false
//	└────────────────────────────┘
//
// `
func (m *moveToSafetyHandler) isObstacleDetected(ev events.UpdateDistanceSensorEvent) bool {
	actualDistance := ev.FrontDistance
	if m.Direction == command.MoveDirectionBackward {
		actualDistance = ev.BackDistance
	}

	if m.obstacleDetected {
		return actualDistance < m.ObstacleTracking.ExitDistance
	}

	return actualDistance < m.ObstacleTracking.EnterDistance
}

func (m *moveToSafetyHandler) isCargoLost(downDistance uint16) bool {
	return downDistance > m.CargoLostDistance
}

func (m *moveToSafetyHandler) stopMovement() {
	if !m.isDriveMotorRunning {
		return
	}

	if err := m.stopDriveMotorFunc(); err != nil {
		m.log.Error("failed to stop motor", slog.Any("error", err))
	}

	m.isDriveMotorRunning = false
}

func (m *moveToSafetyHandler) resumeMovement() {
	if m.isDriveMotorRunning || m.cargoLost {
		return
	}

	if err := m.resumeDriveMotorFunc(); err != nil {
		m.log.Error("failed to resume motor", slog.Any("error", err))
	}

	m.isDriveMotorRunning = true
}
