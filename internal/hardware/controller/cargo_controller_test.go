package controller

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/hardware/picserial"
	"github.com/tbe-team/raybot/internal/logging"
)

func TestCargoController_SetCargoPosition(t *testing.T) {
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

		err := controller.SetCargoPosition(context.Background(), 100, 10)
		require.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":2,
			"data":{
				"target_position":10,
				"max_output":100,
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

		err := controller.SetCargoPosition(context.Background(), 100, 10)
		require.ErrorIs(t, err, ErrCommandACKTimeout)
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

		err := controller.SetCargoPosition(ctx, 100, 10)
		require.ErrorIs(t, err, context.Canceled)
	})
}

func TestCargoController_StopLiftCargoMotor(t *testing.T) {
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

		err := controller.StopLiftCargoMotor(context.Background())
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":2,
			"data":{
				"enable":0,
				"max_output":0,
				"target_position":0
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

		err := controller.StopLiftCargoMotor(context.Background())
		require.ErrorIs(t, err, ErrCommandACKTimeout)
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

		err := controller.StopLiftCargoMotor(ctx)
		require.ErrorIs(t, err, context.Canceled)
	})
}
