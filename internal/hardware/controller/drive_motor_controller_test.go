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

func TestDriveMotorController_MoveForward(t *testing.T) {
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

		err := controller.MoveForward(context.Background(), 10)
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":3,
			"data":{
				"direction":0,
				"speed":10,
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

		err := controller.MoveForward(context.Background(), 10)
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

		err := controller.MoveForward(ctx, 10)
		assert.ErrorIs(t, err, context.Canceled)
	})
}

func TestDriveMotorController_MoveBackward(t *testing.T) {
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

		err := controller.MoveBackward(context.Background(), 10)
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":3,
			"data":{
				"direction":1,
				"speed":10,
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

		err := controller.MoveBackward(context.Background(), 10)
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

		err := controller.MoveBackward(ctx, 10)
		assert.ErrorIs(t, err, context.Canceled)
	})
}

func TestDriveMotorController_StopDriveMotor(t *testing.T) {
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

		err := controller.StopDriveMotor(context.Background(), true)
		assert.NoError(t, err)

		actual := removeMarkers(mockPort.WriteBuffer.Bytes())
		expected := `
		{
			"id":"abc",
			"type":3,
			"data":{
				"direction":0,
				"speed":0,
				"enable":0
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

		err := controller.StopDriveMotor(context.Background(), true)
		assert.ErrorIs(t, err, ErrCommandACKTimeout)
	})

	t.Run("Should fail due to context canceled", func(t *testing.T) {
		log := logging.NewNoopLogger()
		mockPort := &picserial.FakeSerialPort{}
		eventBus := &fakeEventBus{}
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

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := controller.StopDriveMotor(ctx, true)
		assert.ErrorIs(t, err, context.Canceled)
	})
}
