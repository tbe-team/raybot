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

func newScanLocationExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	driveMotorService drivemotor.Service,
	commandRepository command.Repository,
) *commandExecutor[command.ScanLocationInputs, command.ScanLocationOutputs] {
	handler := scanLocationHandler{
		log:               log,
		subscriber:        subscriber,
		driveMotorService: driveMotorService,
	}

	return newCommandExecutor(
		handler.Handle,
		Hooks[command.ScanLocationOutputs]{
			OnCancel: func(ctx context.Context) {
				if err := driveMotorService.Stop(ctx); err != nil {
					log.Error("failed to stop drive motor", slog.Any("error", err))
				}
			},
		},
		log,
		commandRepository,
	)
}

type scanLocationHandler struct {
	log               *slog.Logger
	subscriber        eventbus.Subscriber
	driveMotorService drivemotor.Service
}

func (h scanLocationHandler) Handle(ctx context.Context, _ command.ScanLocationInputs) (command.ScanLocationOutputs, error) {
	outputs := command.ScanLocationOutputs{
		Locations: []command.Location{},
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		outputs.Locations = h.recordLocationsUntilLoopedBack(ctx)
	}()

	if err := h.driveMotorService.MoveForward(ctx, drivemotor.MoveForwardParams{
		Speed: 100,
	}); err != nil {
		return outputs, fmt.Errorf("failed to move forward: %w", err)
	}

	// wait for recording to finish
	wg.Wait()

	if err := h.driveMotorService.Stop(ctx); err != nil {
		return outputs, fmt.Errorf("failed to stop drive motor: %w", err)
	}

	return outputs, nil
}

func (h scanLocationHandler) recordLocationsUntilLoopedBack(ctx context.Context) []command.Location {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		h.log.Debug("stop recording location")
		cancel()
	}()

	locs := []command.Location{}

	doneCh := make(chan struct{})
	h.log.Debug("start recording location")
	h.subscriber.Subscribe(ctx, events.LocationUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateLocationEvent)
		if !ok {
			h.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		// if the first location is reached again, close the channel
		if len(locs) > 0 && ev.Location == locs[0].Location {
			close(doneCh)
			return
		}

		h.log.Debug("record location", slog.String("location", ev.Location))
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
