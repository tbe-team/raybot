package executor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/drivemotor"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type scanLocationExecutor struct {
	log               *slog.Logger
	subscriber        eventbus.Subscriber
	driveMotorService drivemotor.Service
}

func newScanLocationExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	driveMotorService drivemotor.Service,
) CommandExecutor[command.ScanLocationInputs, command.ScanLocationOutputs] {
	return scanLocationExecutor{
		log:               log,
		subscriber:        subscriber,
		driveMotorService: driveMotorService,
	}
}

func (e scanLocationExecutor) Execute(ctx context.Context, _ command.ScanLocationInputs) (command.ScanLocationOutputs, error) {
	outputs := command.ScanLocationOutputs{
		Locations: []command.Location{},
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		outputs.Locations = e.recordLocationsUntilLoopedBack(ctx)
	}()

	if err := e.driveMotorService.MoveForward(ctx, drivemotor.MoveForwardParams{
		Speed: 100,
	}); err != nil {
		return outputs, fmt.Errorf("failed to move forward: %w", err)
	}

	// wait for recording to finish
	wg.Wait()

	if err := e.driveMotorService.Stop(ctx); err != nil {
		return outputs, fmt.Errorf("failed to stop drive motor: %w", err)
	}

	return outputs, nil
}

func (e scanLocationExecutor) OnCancel(ctx context.Context) error {
	if err := e.driveMotorService.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop drive motor: %w", err)
	}
	return nil
}

func (e scanLocationExecutor) recordLocationsUntilLoopedBack(ctx context.Context) []command.Location {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		e.log.Debug("stop recording location")
		cancel()
	}()

	locs := []command.Location{}

	doneCh := make(chan struct{})
	e.log.Debug("start recording location")
	e.subscriber.Subscribe(ctx, events.LocationUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateLocationEvent)
		if !ok {
			e.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		// if the first location is reached again, close the channel
		if len(locs) > 0 && ev.Location == locs[0].Location {
			close(doneCh)
			return
		}

		e.log.Debug("record location", slog.String("location", ev.Location))
		locs = append(locs, command.Location{
			Location:  ev.Location,
			ScannedAt: time.Now(),
		})
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}

	return locs
}
