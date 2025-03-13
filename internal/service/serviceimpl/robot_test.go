package serviceimpl

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository/mocks"
	"github.com/tbe-team/raybot/internal/service"
	dbmocks "github.com/tbe-team/raybot/internal/storage/db/mocks"
	"github.com/tbe-team/raybot/pkg/validator"
)

func TestRobotService(t *testing.T) {
	validator := validator.New()
	ctx := context.Background()

	t.Run("test GetRobotState", func(t *testing.T) {
		tests := []struct {
			name          string
			mock          func(_ *mocks.FakeRobotStateRepository, _ *dbmocks.FakeProvider)
			expectedState model.RobotState
			expectedError bool
		}{
			{
				name: "successful get robot state",
				mock: func(robotStateRepo *mocks.FakeRobotStateRepository, dbProvider *dbmocks.FakeProvider) {
					dbProvider.EXPECT().DB().Return(nil)
					robotStateRepo.EXPECT().GetRobotState(ctx, mock.Anything).Return(model.RobotState{
						Battery: model.BatteryState{
							Current: 100,
						},
					}, nil)
				},
				expectedState: model.RobotState{
					Battery: model.BatteryState{
						Current: 100,
					},
				},
				expectedError: false,
			},
			{
				name: "get robot state failed",
				mock: func(robotStateRepo *mocks.FakeRobotStateRepository, dbProvider *dbmocks.FakeProvider) {
					dbProvider.EXPECT().DB().Return(nil)
					robotStateRepo.EXPECT().GetRobotState(ctx, mock.Anything).Return(model.RobotState{}, assert.AnError)
				},
				expectedState: model.RobotState{},
				expectedError: true,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				robotStateRepo := mocks.NewFakeRobotStateRepository(t)
				dbProvider := dbmocks.NewFakeProvider(t)
				s := NewRobotService(robotStateRepo, dbProvider, validator)

				tc.mock(robotStateRepo, dbProvider)

				state, err := s.GetRobotState(ctx)
				if tc.expectedError {
					assert.Error(t, err)
					return
				}
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedState, state)
			})
		}
	})

	t.Run("test UpdateRobotState", func(t *testing.T) {
		tests := []struct {
			name          string
			params        service.UpdateRobotStateParams
			mock          func(_ *mocks.FakeRobotStateRepository, _ *dbmocks.FakeProvider)
			expectedState model.RobotState
			expectedError bool
		}{
			{
				name: "validation failed",
				params: service.UpdateRobotStateParams{
					SetDriveMotor: true,
					DriveMotor: service.DriveMotorParams{
						Direction: 99, // Invalid direction
					},
				},
				mock: func(_ *mocks.FakeRobotStateRepository, _ *dbmocks.FakeProvider) {
				},
				expectedState: model.RobotState{},
				expectedError: true,
			},
			{
				name: "get robot state failed",
				params: service.UpdateRobotStateParams{
					SetBattery: true,
					Battery: service.BatteryParams{
						Current: 100,
					},
				},
				mock: func(robotStateRepo *mocks.FakeRobotStateRepository, dbProvider *dbmocks.FakeProvider) {
					dbProvider.EXPECT().DB().Return(nil)
					robotStateRepo.EXPECT().GetRobotState(ctx, mock.Anything).Return(model.RobotState{}, assert.AnError)
				},
				expectedState: model.RobotState{},
				expectedError: true,
			},
			{
				name: "update robot state failed",
				params: service.UpdateRobotStateParams{
					SetBattery: true,
					Battery: service.BatteryParams{
						Current: 100,
					},
				},
				mock: func(robotStateRepo *mocks.FakeRobotStateRepository, dbProvider *dbmocks.FakeProvider) {
					dbProvider.EXPECT().DB().Return(nil)
					robotStateRepo.EXPECT().GetRobotState(ctx, mock.Anything).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything, mock.Anything).Return(assert.AnError)
				},
				expectedState: model.RobotState{},
				expectedError: true,
			},
			{
				name: "update battery state successful",
				params: service.UpdateRobotStateParams{
					SetBattery: true,
					Battery: service.BatteryParams{
						Current:      100,
						Temp:         25,
						Voltage:      12000,
						CellVoltages: []uint16{4000, 4000, 4000},
						Percent:      80,
						Fault:        0,
						Health:       100,
					},
				},
				mock: func(robotStateRepo *mocks.FakeRobotStateRepository, dbProvider *dbmocks.FakeProvider) {
					dbProvider.EXPECT().DB().Return(nil)
					robotStateRepo.EXPECT().GetRobotState(ctx, mock.Anything).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything, mock.Anything).Return(nil)
				},
				expectedState: model.RobotState{},
				expectedError: false,
			},
			{
				name: "update charge state successful",
				params: service.UpdateRobotStateParams{
					SetCharge: true,
					Charge: service.ChargeParams{
						CurrentLimit: 1000,
						Enabled:      true,
					},
				},
				mock: func(robotStateRepo *mocks.FakeRobotStateRepository, dbProvider *dbmocks.FakeProvider) {
					dbProvider.EXPECT().DB().Return(nil)
					robotStateRepo.EXPECT().GetRobotState(ctx, mock.Anything).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything, mock.Anything).Return(nil)
				},
				expectedState: model.RobotState{},
				expectedError: false,
			},
			{
				name: "update discharge state successful",
				params: service.UpdateRobotStateParams{
					SetDischarge: true,
					Discharge: service.DischargeParams{
						CurrentLimit: 500,
						Enabled:      true,
					},
				},
				mock: func(robotStateRepo *mocks.FakeRobotStateRepository, dbProvider *dbmocks.FakeProvider) {
					dbProvider.EXPECT().DB().Return(nil)
					robotStateRepo.EXPECT().GetRobotState(ctx, mock.Anything).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything, mock.Anything).Return(nil)
				},
				expectedState: model.RobotState{},
				expectedError: false,
			},
			{
				name: "update distance sensor state successful",
				params: service.UpdateRobotStateParams{
					SetDistanceSensor: true,
					DistanceSensor: service.DistanceSensorParams{
						FrontDistance: 100,
						BackDistance:  200,
						DownDistance:  50,
					},
				},
				mock: func(robotStateRepo *mocks.FakeRobotStateRepository, dbProvider *dbmocks.FakeProvider) {
					dbProvider.EXPECT().DB().Return(nil)
					robotStateRepo.EXPECT().GetRobotState(ctx, mock.Anything).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything, mock.Anything).Return(nil)
				},
				expectedState: model.RobotState{},
				expectedError: false,
			},
			{
				name: "update lift motor state successful",
				params: service.UpdateRobotStateParams{
					SetLiftMotor: true,
					LiftMotor: service.LiftMotorParams{
						CurrentPosition: 50,
						TargetPosition:  100,
						IsRunning:       true,
						Enabled:         true,
					},
				},
				mock: func(robotStateRepo *mocks.FakeRobotStateRepository, dbProvider *dbmocks.FakeProvider) {
					dbProvider.EXPECT().DB().Return(nil)
					robotStateRepo.EXPECT().GetRobotState(ctx, mock.Anything).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything, mock.Anything).Return(nil)
				},
				expectedState: model.RobotState{},
				expectedError: false,
			},
			{
				name: "update drive motor state successful",
				params: service.UpdateRobotStateParams{
					SetDriveMotor: true,
					DriveMotor: service.DriveMotorParams{
						Direction: model.DriveMotorDirectionForward,
						Speed:     75,
						IsRunning: true,
						Enabled:   true,
					},
				},
				mock: func(robotStateRepo *mocks.FakeRobotStateRepository, dbProvider *dbmocks.FakeProvider) {
					dbProvider.EXPECT().DB().Return(nil)
					robotStateRepo.EXPECT().GetRobotState(ctx, mock.Anything).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything, mock.Anything).Return(nil)
				},
				expectedState: model.RobotState{},
				expectedError: false,
			},
			{
				name: "update multiple states successful",
				params: service.UpdateRobotStateParams{
					SetBattery: true,
					Battery: service.BatteryParams{
						Current: 100,
						Voltage: 12000,
					},
					SetDriveMotor: true,
					DriveMotor: service.DriveMotorParams{
						Direction: model.DriveMotorDirectionBackward,
						Speed:     50,
					},
				},
				mock: func(robotStateRepo *mocks.FakeRobotStateRepository, dbProvider *dbmocks.FakeProvider) {
					dbProvider.EXPECT().DB().Return(nil)
					robotStateRepo.EXPECT().GetRobotState(ctx, mock.Anything).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything, mock.Anything).Return(nil)
				},
				expectedState: model.RobotState{},
				expectedError: false,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				robotStateRepo := mocks.NewFakeRobotStateRepository(t)
				dbProvider := dbmocks.NewFakeProvider(t)
				s := NewRobotService(robotStateRepo, dbProvider, validator)

				tc.mock(robotStateRepo, dbProvider)

				state, err := s.UpdateRobotState(ctx, tc.params)
				if tc.expectedError {
					assert.Error(t, err)
					return
				}

				assert.NoError(t, err)

				// Check that updatedAt fields were set properly
				if tc.params.SetBattery {
					assert.WithinDuration(t, time.Now(), state.Battery.UpdatedAt, 2*time.Second)
				}
				if tc.params.SetCharge {
					assert.WithinDuration(t, time.Now(), state.Charge.UpdatedAt, 2*time.Second)
				}
				if tc.params.SetDischarge {
					assert.WithinDuration(t, time.Now(), state.Discharge.UpdatedAt, 2*time.Second)
				}
				if tc.params.SetDistanceSensor {
					assert.WithinDuration(t, time.Now(), state.DistanceSensor.UpdatedAt, 2*time.Second)
				}
				if tc.params.SetLiftMotor {
					assert.WithinDuration(t, time.Now(), state.LiftMotor.UpdatedAt, 2*time.Second)
				}
				if tc.params.SetDriveMotor {
					assert.WithinDuration(t, time.Now(), state.DriveMotor.UpdatedAt, 2*time.Second)
				}
			})
		}
	})
}
