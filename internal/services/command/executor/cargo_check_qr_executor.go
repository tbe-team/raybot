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
	handler := cargoCheckQRHandler{
		log:        log,
		subscriber: subscriber,
	}

	return newCommandExecutor(
		handler.Handle,
		Hooks[command.CargoCheckQROutputs]{},
		log,
		commandRepository,
	)
}

type cargoCheckQRHandler struct {
	log        *slog.Logger
	subscriber eventbus.Subscriber
}

func (c cargoCheckQRHandler) Handle(ctx context.Context, inputs command.CargoCheckQRInputs) (command.CargoCheckQROutputs, error) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.trackingCargoQRCodeUntilMatched(ctx, inputs.QRCode)
	}()

	// wait for tracking to finish
	wg.Wait()

	return command.CargoCheckQROutputs{}, nil
}

func (c cargoCheckQRHandler) trackingCargoQRCodeUntilMatched(ctx context.Context, qrCode string) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		c.log.Debug("stop tracking cargo qr code")
		cancel()
	}()

	doneCh := make(chan struct{})
	c.log.Debug("start tracking cargo qr code", slog.Any("qr_code", qrCode))
	c.subscriber.Subscribe(ctx, events.CargoQRCodeUpdatedTopic, func(_ context.Context, msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.CargoQRCodeUpdatedEvent)
		if !ok {
			c.log.Error("invalid event", slog.Any("event", msg.Payload))
			return
		}

		if ev.QRCode == qrCode {
			c.log.Debug("cargo qr code matched", slog.Any("qrcode", ev.QRCode))
			close(doneCh)
		}
	})

	select {
	case <-doneCh:
	case <-ctx.Done():
	}
}
