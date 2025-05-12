package emergencyimpl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	commandmocks "github.com/tbe-team/raybot/internal/services/command/mocks"
	drivemotormocks "github.com/tbe-team/raybot/internal/services/drivemotor/mocks"
	"github.com/tbe-team/raybot/internal/services/emergency"
	liftmotormocks "github.com/tbe-team/raybot/internal/services/liftmotor/mocks"
)

func TestService(t *testing.T) {
	t.Run("Get emergency state", func(t *testing.T) {
		t.Run("Should return the emergency state successfully", func(t *testing.T) {
			commandService := commandmocks.NewFakeService(t)
			driveMotorService := drivemotormocks.NewFakeService(t)
			liftMotorService := liftmotormocks.NewFakeService(t)
			emergencyStateRepository := NewEmergencyStateRepository()

			err := emergencyStateRepository.UpdateEmergencyState(context.Background(), emergency.State{Locked: true})
			assert.NoError(t, err)

			service := NewService(commandService, driveMotorService, liftMotorService, emergencyStateRepository)

			state, err := service.GetEmergencyState(context.Background())
			assert.NoError(t, err)
			assert.True(t, state.Locked)
		})
	})

	t.Run("Stop operation", func(t *testing.T) {
		t.Run("Should stop the command processing, drive motor and lift motor successfully", func(t *testing.T) {
			commandService := commandmocks.NewFakeService(t)
			driveMotorService := drivemotormocks.NewFakeService(t)
			liftMotorService := liftmotormocks.NewFakeService(t)
			emergencyStateRepository := NewEmergencyStateRepository()

			commandService.EXPECT().LockProcessingCommand(context.Background()).Return(nil)
			driveMotorService.EXPECT().Stop(context.Background()).Return(nil)
			liftMotorService.EXPECT().Stop(context.Background()).Return(nil)

			service := NewService(commandService, driveMotorService, liftMotorService, emergencyStateRepository)

			err := service.StopOperation(context.Background())
			assert.NoError(t, err)

			state, err := emergencyStateRepository.GetEmergencyState(context.Background())
			assert.NoError(t, err)
			assert.True(t, state.Locked)
		})

		t.Run("Should return an error if lock processing command fails", func(t *testing.T) {
			commandService := commandmocks.NewFakeService(t)
			driveMotorService := drivemotormocks.NewFakeService(t)
			liftMotorService := liftmotormocks.NewFakeService(t)
			emergencyStateRepository := NewEmergencyStateRepository()

			commandService.EXPECT().LockProcessingCommand(context.Background()).Return(assert.AnError)
			driveMotorService.AssertNotCalled(t, "Stop")
			liftMotorService.AssertNotCalled(t, "Stop")

			service := NewService(commandService, driveMotorService, liftMotorService, emergencyStateRepository)

			err := service.StopOperation(context.Background())
			assert.ErrorIs(t, err, assert.AnError)

			state, err := emergencyStateRepository.GetEmergencyState(context.Background())
			assert.NoError(t, err)
			assert.False(t, state.Locked)
		})

		t.Run("Should return an error if the drive motor service fails", func(t *testing.T) {
			commandService := commandmocks.NewFakeService(t)
			driveMotorService := drivemotormocks.NewFakeService(t)
			liftMotorService := liftmotormocks.NewFakeService(t)
			emergencyStateRepository := NewEmergencyStateRepository()

			commandService.EXPECT().LockProcessingCommand(context.Background()).Return(nil)
			driveMotorService.EXPECT().Stop(context.Background()).Return(assert.AnError)
			liftMotorService.AssertNotCalled(t, "Stop")

			service := NewService(commandService, driveMotorService, liftMotorService, emergencyStateRepository)

			err := service.StopOperation(context.Background())
			assert.ErrorIs(t, err, assert.AnError)

			state, err := emergencyStateRepository.GetEmergencyState(context.Background())
			assert.NoError(t, err)
			assert.False(t, state.Locked)
		})

		t.Run("Should return an error if the lift motor service fails", func(t *testing.T) {
			commandService := commandmocks.NewFakeService(t)
			driveMotorService := drivemotormocks.NewFakeService(t)
			liftMotorService := liftmotormocks.NewFakeService(t)
			emergencyStateRepository := NewEmergencyStateRepository()

			commandService.EXPECT().LockProcessingCommand(context.Background()).Return(nil)
			driveMotorService.EXPECT().Stop(context.Background()).Return(nil)
			liftMotorService.EXPECT().Stop(context.Background()).Return(assert.AnError)

			service := NewService(commandService, driveMotorService, liftMotorService, emergencyStateRepository)

			err := service.StopOperation(context.Background())
			assert.ErrorIs(t, err, assert.AnError)

			state, err := emergencyStateRepository.GetEmergencyState(context.Background())
			assert.NoError(t, err)
			assert.False(t, state.Locked)
		})
	})

	t.Run("Resume operation", func(t *testing.T) {
		t.Run("Should unlock the command processing successfully", func(t *testing.T) {
			commandService := commandmocks.NewFakeService(t)
			driveMotorService := drivemotormocks.NewFakeService(t)
			liftMotorService := liftmotormocks.NewFakeService(t)
			emergencyStateRepository := NewEmergencyStateRepository()

			commandService.EXPECT().UnlockProcessingCommand(context.Background()).Return(nil)

			service := NewService(commandService, driveMotorService, liftMotorService, emergencyStateRepository)

			err := service.ResumeOperation(context.Background())
			assert.NoError(t, err)

			state, err := emergencyStateRepository.GetEmergencyState(context.Background())
			assert.NoError(t, err)
			assert.False(t, state.Locked)
		})

		t.Run("Should return an error if the command processing lock fails", func(t *testing.T) {
			commandService := commandmocks.NewFakeService(t)
			driveMotorService := drivemotormocks.NewFakeService(t)
			liftMotorService := liftmotormocks.NewFakeService(t)
			emergencyStateRepository := NewEmergencyStateRepository()

			commandService.EXPECT().UnlockProcessingCommand(context.Background()).Return(assert.AnError)

			service := NewService(commandService, driveMotorService, liftMotorService, emergencyStateRepository)

			err := service.ResumeOperation(context.Background())
			assert.ErrorIs(t, err, assert.AnError)

			state, err := emergencyStateRepository.GetEmergencyState(context.Background())
			assert.NoError(t, err)
			assert.False(t, state.Locked)
		})
	})
}
