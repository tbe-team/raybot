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
	outputs := command.ScanLocationOutputs{
		Locations: []command.Location{},
	}

	return newCommandExecutor(
		func(ctx context.Context, _ command.ScanLocationInputs) (command.ScanLocationOutputs, error) {
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				outputs.Locations = recordLocationsUntilLoopedBack(ctx, log, subscriber)
			}()

			if err := driveMotorService.MoveForward(ctx, drivemotor.MoveForwardParams{
				Speed: 100,
			}); err != nil {
				return outputs, fmt.Errorf("failed to move forward: %w", err)
			}

			wg.Wait()

			if err := driveMotorService.Stop(ctx); err != nil {
				return outputs, fmt.Errorf("failed to stop drive motor: %w", err)
			}

			return outputs, nil
		},
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

func recordLocationsUntilLoopedBack(
	ctx context.Context,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
) []command.Location {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		log.Debug("stop recording location")
		cancel()
	}()

	locs := []command.Location{}

	doneCh := make(chan struct{})
	log.Debug("start recording location")
	subscriber.Subscribe(ctx, events.LocationUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.UpdateLocationEvent)
		if !ok {
			log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		// if the first location is reached again, close the channel
		if len(locs) > 0 && ev.Location == locs[0].Location {
			close(doneCh)
			return
		}

		log.Debug("record location", slog.String("location", ev.Location))
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
