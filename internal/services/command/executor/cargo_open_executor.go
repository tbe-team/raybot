package executor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/cargo"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

func newCargoOpenExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	cargoService cargo.Service,
	commandRepository command.Repository,
) *commandExecutor[command.CargoOpenInputs] {
	return newCommandExecutor(func(ctx context.Context, _ command.CargoOpenInputs) error {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			trackingCargoDoorUntilOpen(ctx, log, subscriber)
		}()

		if err := cargoService.OpenCargoDoor(ctx); err != nil {
			return fmt.Errorf("failed to open cargo: %w", err)
		}

		// wait for cargo door open to finish
		wg.Wait()

		return nil
	},
		Hooks{},
		log,
		commandRepository,
	)
}

func trackingCargoDoorUntilOpen(
	ctx context.Context,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		log.Debug("stop tracking cargo door open")
		cancel()
	}()

	doneCh := make(chan struct{})
	log.Debug("start tracking cargo door open")
	subscriber.Subscribe(ctx, events.CargoDoorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.CargoDoorUpdatedEvent)
		if !ok {
			log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if ev.IsOpen {
			log.Debug("cargo door open completed")
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
