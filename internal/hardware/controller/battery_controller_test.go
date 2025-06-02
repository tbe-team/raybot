package controller

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/hardware/picserial"
	"github.com/tbe-team/raybot/internal/logging"
)

func TestBatteryController_ConfigBatteryCharge(t *testing.T) {
	t.Run("Should success", func(t *testing.T) {
		log := logging.NewNoopLogger()
		mockPort := &picserial.FakeSerialPort{}
		eventBus := &fakeEventBus{
			expectedPayload: events.PICCmdAckEvent{
				ID:      "abc",
				Success: true,
			},
		}
		controller := controller{
			cfg: config.Hardware{
				PIC: config.PIC{
					CommandACKTimeout: 10 * time.Millisecond,
				},
			},
			log:             log,
			subscriber:      eventBus,
			picSerialClient: picserial.NewClientWithPort(mockPort),
			genIDFunc:       func() string { return "abc" },
		}

		err := controller.ConfigBatteryCharge(context.Background(), 10, true)
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
	{
		"id":"abc",
		"type":0,
		"data":{
			"current_limit":10,
			"enable":1
		}
	}`
		assert.JSONEq(t, expected, string(actual))
	})

	t.Run("Should fail due to timeout", func(t *testing.T) {
		log := logging.NewNoopLogger()
		mockPort := &picserial.FakeSerialPort{}
		eventBus := &fakeEventBus{}
		controller := controller{
			cfg: config.Hardware{
				PIC: config.PIC{
					EnableACK: true,
				},
			},
			log:             log,
			subscriber:      eventBus,
			picSerialClient: picserial.NewClientWithPort(mockPort),
			genIDFunc:       func() string { return "abc" },
		}

		err := controller.ConfigBatteryCharge(context.Background(), 10, true)
		assert.ErrorIs(t, err, ErrCommandACKTimeout)
	})

	t.Run("Should fail due to context canceled", func(t *testing.T) {
		log := logging.NewNoopLogger()
		mockPort := &picserial.FakeSerialPort{}
		eventBus := &fakeEventBus{}
		controller := controller{
			cfg:             config.Hardware{},
			log:             log,
			subscriber:      eventBus,
			picSerialClient: picserial.NewClientWithPort(mockPort),
			genIDFunc:       func() string { return "abc" },
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := controller.ConfigBatteryCharge(ctx, 10, true)
		assert.ErrorIs(t, err, context.Canceled)
	})
}

func TestBatteryController_ConfigBatteryDischarge(t *testing.T) {
	t.Run("Should success", func(t *testing.T) {
		log := logging.NewNoopLogger()
		mockPort := &picserial.FakeSerialPort{}
		eventBus := &fakeEventBus{
			expectedPayload: events.PICCmdAckEvent{
				ID:      "abc",
				Success: true,
			},
		}
		controller := controller{
			cfg: config.Hardware{
				PIC: config.PIC{
					CommandACKTimeout: 10 * time.Millisecond,
				},
			},
			log:             log,
			subscriber:      eventBus,
			picSerialClient: picserial.NewClientWithPort(mockPort),
			genIDFunc:       func() string { return "abc" },
		}

		err := controller.ConfigBatteryDischarge(context.Background(), 10, true)
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":1,
			"data":{
				"current_limit":10,
				"enable":1
			}
		}`
		assert.JSONEq(t, expected, string(actual))
	})

	t.Run("Should fail due to timeout", func(t *testing.T) {
		log := logging.NewNoopLogger()
		mockPort := &picserial.FakeSerialPort{}
		eventBus := &fakeEventBus{}
		controller := controller{
			cfg: config.Hardware{
				PIC: config.PIC{
					EnableACK: true,
				},
			},
			log:             log,
			subscriber:      eventBus,
			picSerialClient: picserial.NewClientWithPort(mockPort),
			genIDFunc:       func() string { return "abc" },
		}

		err := controller.ConfigBatteryDischarge(context.Background(), 10, true)
		assert.ErrorIs(t, err, ErrCommandACKTimeout)
	})

	t.Run("Should fail due to context canceled", func(t *testing.T) {
		log := logging.NewNoopLogger()
		mockPort := &picserial.FakeSerialPort{}
		eventBus := &fakeEventBus{}
		controller := controller{
			cfg:             config.Hardware{},
			log:             log,
			subscriber:      eventBus,
			picSerialClient: picserial.NewClientWithPort(mockPort),
			genIDFunc:       func() string { return "abc" },
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := controller.ConfigBatteryDischarge(ctx, 10, true)
		assert.ErrorIs(t, err, context.Canceled)
	})
}
