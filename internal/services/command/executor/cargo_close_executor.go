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

func newCargoCloseExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	cargoService cargo.Service,
	commandRepository command.Repository,
) *commandExecutor[command.CargoCloseInputs, command.CargoCloseOutputs] {
	return newCommandExecutor(
		func(ctx context.Context, _ command.CargoCloseInputs) (command.CargoCloseOutputs, error) {
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				trackingCargoDoorUntilClosed(ctx, log, subscriber)
			}()

			if err := cargoService.CloseCargoDoor(ctx); err != nil {
				return command.CargoCloseOutputs{}, fmt.Errorf("failed to close cargo door: %w", err)
			}

			// wait for cargo door close to finish
			wg.Wait()

			return command.CargoCloseOutputs{}, nil
		},
		Hooks[command.CargoCloseOutputs]{},
		log,
		commandRepository,
	)
}

func trackingCargoDoorUntilClosed(
	ctx context.Context,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		log.Debug("stop tracking cargo door")
		cancel()
	}()

	doneCh := make(chan struct{})
	log.Debug("start tracking cargo door")
	subscriber.Subscribe(ctx, events.CargoDoorUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.CargoDoorUpdatedEvent)
		if !ok {
			log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if !ev.IsOpen {
			log.Debug("cargo door closed")
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
