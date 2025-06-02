package controller

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/hardware/espserial"
	"github.com/tbe-team/raybot/internal/logging"
)

func TestCargoDoorController_OpenCargoDoor(t *testing.T) {
	t.Run("Should success", func(t *testing.T) {
		log := logging.NewNoopLogger()
		mockPort := &espserial.FakeSerialPort{}
		eventBus := &fakeEventBus{
			expectedPayload: events.ESPCmdAckEvent{
				ID:      "abc",
				Success: true,
			},
		}
		controller := controller{
			cfg: config.Hardware{
				ESP: config.ESP{
					CommandACKTimeout: 10 * time.Millisecond,
				},
			},
			log:             log,
			subscriber:      eventBus,
			espSerialClient: espserial.NewClientWithPort(mockPort),
			genIDFunc:       func() string { return "abc" },
		}

		err := controller.OpenCargoDoor(context.Background(), 10)
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":0,
			"data":{
				"state":1,
				"speed":10,
				"enable":1
			}
		}`
		assert.JSONEq(t, expected, string(actual))
	})

	t.Run("Should fail due to timeout", func(t *testing.T) {
		log := logging.NewNoopLogger()
		mockPort := &espserial.FakeSerialPort{}
		eventBus := &fakeEventBus{}
		controller := controller{
			cfg: config.Hardware{
				ESP: config.ESP{
					EnableACK: true,
				},
			},
			log:             log,
			subscriber:      eventBus,
			espSerialClient: espserial.NewClientWithPort(mockPort),
			genIDFunc:       func() string { return "abc" },
		}

		err := controller.OpenCargoDoor(context.Background(), 10)
		assert.ErrorIs(t, err, ErrCommandACKTimeout)
	})

	t.Run("Should fail due to context canceled", func(t *testing.T) {
		log := logging.NewNoopLogger()
		mockPort := &espserial.FakeSerialPort{}
		eventBus := &fakeEventBus{}
		controller := controller{
			cfg:             config.Hardware{},
			log:             log,
			subscriber:      eventBus,
			espSerialClient: espserial.NewClientWithPort(mockPort),
			genIDFunc:       func() string { return "abc" },
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := controller.OpenCargoDoor(ctx, 10)
		assert.ErrorIs(t, err, context.Canceled)
	})
}

func TestCargoDoorController_CloseCargoDoor(t *testing.T) {
	t.Run("Should success", func(t *testing.T) {
		log := logging.NewNoopLogger()
		mockPort := &espserial.FakeSerialPort{}
		eventBus := &fakeEventBus{
			expectedPayload: events.ESPCmdAckEvent{
				ID:      "abc",
				Success: true,
			},
		}
		controller := controller{
			cfg: config.Hardware{
				ESP: config.ESP{
					CommandACKTimeout: 10 * time.Millisecond,
				},
			},
			log:             log,
			subscriber:      eventBus,
			espSerialClient: espserial.NewClientWithPort(mockPort),
			genIDFunc:       func() string { return "abc" },
		}

		err := controller.CloseCargoDoor(context.Background(), 10)
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":0,
			"data":{
				"state":0,
				"speed":10,
				"enable":1
			}
		}`
		assert.JSONEq(t, expected, string(actual))
	})

	t.Run("Should fail due to timeout", func(t *testing.T) {
		log := logging.NewNoopLogger()
		mockPort := &espserial.FakeSerialPort{}
		eventBus := &fakeEventBus{}
		controller := controller{
			cfg: config.Hardware{
				ESP: config.ESP{
					EnableACK: true,
				},
			},
			log:             log,
			subscriber:      eventBus,
			espSerialClient: espserial.NewClientWithPort(mockPort),
			genIDFunc:       func() string { return "abc" },
		}

		err := controller.CloseCargoDoor(context.Background(), 10)
		assert.ErrorIs(t, err, ErrCommandACKTimeout)
	})

	t.Run("Should fail due to context canceled", func(t *testing.T) {
		log := logging.NewNoopLogger()
		mockPort := &espserial.FakeSerialPort{}
		eventBus := &fakeEventBus{}
		controller := controller{
			cfg:             config.Hardware{},
			log:             log,
			subscriber:      eventBus,
			espSerialClient: espserial.NewClientWithPort(mockPort),
			genIDFunc:       func() string { return "abc" },
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := controller.CloseCargoDoor(ctx, 10)
		assert.ErrorIs(t, err, context.Canceled)
	})
}
