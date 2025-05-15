package systemimpl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tbe-team/raybot/internal/logging"
	commandmocks "github.com/tbe-team/raybot/internal/services/command/mocks"
	drivemotormocks "github.com/tbe-team/raybot/internal/services/drivemotor/mocks"
	liftmotormocks "github.com/tbe-team/raybot/internal/services/liftmotor/mocks"
)

func TestService(t *testing.T) {
	t.Run("Stop emergency", func(t *testing.T) {
		t.Run("Should stop all motors and commands", func(t *testing.T) {
			log := logging.NewNoopLogger()
			commandService := commandmocks.NewFakeService(t)
			driveMotorService := drivemotormocks.NewFakeService(t)
			liftMotorService := liftmotormocks.NewFakeService(t)

			commandService.EXPECT().CancelAllRunningCommands(context.Background()).Return(nil)
			driveMotorService.EXPECT().Stop(context.Background()).Return(nil)
			liftMotorService.EXPECT().Stop(context.Background()).Return(nil)

			service := NewService(log, commandService, driveMotorService, liftMotorService)

			err := service.StopEmergency(context.Background())
			assert.NoError(t, err)
		})

		t.Run("Should return error if cancel all running commands fails", func(t *testing.T) {
			log := logging.NewNoopLogger()
			commandService := commandmocks.NewFakeService(t)
			driveMotorService := drivemotormocks.NewFakeService(t)
			liftMotorService := liftmotormocks.NewFakeService(t)

			commandService.EXPECT().CancelAllRunningCommands(context.Background()).Return(assert.AnError)
			driveMotorService.AssertNotCalled(t, "Stop")
			liftMotorService.AssertNotCalled(t, "Stop")

			service := NewService(log, commandService, driveMotorService, liftMotorService)

			err := service.StopEmergency(context.Background())
			assert.Error(t, err)
		})

		t.Run("Should return error if stop drive motor fails", func(t *testing.T) {
			log := logging.NewNoopLogger()
			commandService := commandmocks.NewFakeService(t)
			driveMotorService := drivemotormocks.NewFakeService(t)
			liftMotorService := liftmotormocks.NewFakeService(t)

			commandService.EXPECT().CancelAllRunningCommands(context.Background()).Return(nil)
			driveMotorService.EXPECT().Stop(context.Background()).Return(assert.AnError)
			liftMotorService.AssertNotCalled(t, "Stop")

			service := NewService(log, commandService, driveMotorService, liftMotorService)

			err := service.StopEmergency(context.Background())
			assert.Error(t, err)
		})

		t.Run("Should return error if stop lift motor fails", func(t *testing.T) {
			log := logging.NewNoopLogger()
			commandService := commandmocks.NewFakeService(t)
			driveMotorService := drivemotormocks.NewFakeService(t)
			liftMotorService := liftmotormocks.NewFakeService(t)

			commandService.EXPECT().CancelAllRunningCommands(context.Background()).Return(nil)
			driveMotorService.EXPECT().Stop(context.Background()).Return(nil)
			liftMotorService.EXPECT().Stop(context.Background()).Return(assert.AnError)

			service := NewService(log, commandService, driveMotorService, liftMotorService)

			err := service.StopEmergency(context.Background())
			assert.Error(t, err)
		})
	})
}
