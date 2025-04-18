package executor

import (
	"context"
	"log/slog"
	"sync"

	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

func newCargoCheckQRExecutor(
	log *slog.Logger,
	subscriber eventbus.Subscriber,
	commandRepository command.Repository,
) *commandExecutor[command.CargoCheckQRInputs, command.CargoCheckQROutputs] {
	return newCommandExecutor(
		func(ctx context.Context, inputs command.CargoCheckQRInputs) (command.CargoCheckQROutputs, error) {
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				trackingCargoQRCodeUntilMatched(ctx, inputs.QRCode, log, subscriber)
			}()

			// wait for cargo qr code to be matched
			wg.Wait()

			return command.CargoCheckQROutputs{}, nil
		},
		Hooks[command.CargoCheckQROutputs]{},
		log,
		commandRepository,
	)
}

func trackingCargoQRCodeUntilMatched(
	ctx context.Context,
	qrCode string,
	log *slog.Logger,
	subscriber eventbus.Subscriber,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		log.Debug("stop tracking cargo qr code")
		cancel()
	}()

	doneCh := make(chan struct{})
	log.Debug("start tracking cargo qr code", slog.Any("qr_code", qrCode))
	subscriber.Subscribe(ctx, events.CargoQRCodeUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.CargoQRCodeUpdatedEvent)
		if !ok {
			log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if ev.QRCode == qrCode {
			log.Debug("cargo qr code matched", slog.Any("qrcode", ev.QRCode))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
