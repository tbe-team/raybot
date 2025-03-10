package serviceimpl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository/mocks"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/pkg/validator"
)

func TestPICService(t *testing.T) {
	validator := validator.New()
	ctx := context.Background()
	t.Run("test process serial command ACK", func(t *testing.T) {
		tests := []struct {
			name          string
			params        service.ProcessSerialCommandACK
			mock          func(_ *mocks.FakePICSerialCommandRepository, _ *mocks.FakeRobotStateRepository)
			expectedError bool
		}{
			{
				name: "validation failed",
				params: service.ProcessSerialCommandACK{
					ID:      "",
					Success: true,
				},
				mock: func(_ *mocks.FakePICSerialCommandRepository, _ *mocks.FakeRobotStateRepository) {
				},
				expectedError: true,
			},
			{
				name: "command execution failed",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: false,
				},
				mock: func(_ *mocks.FakePICSerialCommandRepository, _ *mocks.FakeRobotStateRepository) {
				},
				expectedError: true,
			},
			{
				name: "get pic serial command failed",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, _ *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{}, assert.AnError)
				},
				expectedError: true,
			},
			{
				name: "get robot state failed",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, robotStateRepo *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{}, nil)
					robotStateRepo.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, assert.AnError)
				},
				expectedError: true,
			},
			{
				name: "update robot state by command type - battery charge",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, robotStateRepo *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{
						Type: model.PICSerialCommandTypeBatteryCharge,
						Data: model.PICSerialCommandBatteryChargeData{
							CurrentLimit: 1,
							Enable:       true,
						},
					}, nil)
					robotStateRepo.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything).Return(nil)
					picCommandSerialRepo.EXPECT().DeletePICSerialCommand(ctx, "123").Return(nil)
				},
				expectedError: false,
			},
			{
				name: "update robot state by command type - battery discharge",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, robotStateRepo *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{
						Type: model.PICSerialCommandTypeBatteryDischarge,
						Data: model.PICSerialCommandBatteryDischargeData{
							CurrentLimit: 2,
							Enable:       true,
						},
					}, nil)
					robotStateRepo.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything).Return(nil)
					picCommandSerialRepo.EXPECT().DeletePICSerialCommand(ctx, "123").Return(nil)
				},
				expectedError: false,
			},
			{
				name: "update robot state by command type - lift motor",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, robotStateRepo *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{
						Type: model.PICSerialCommandTypeLiftMotor,
						Data: model.PICSerialCommandBatteryLiftMotorData{
							TargetPosition: 100,
							Enable:         true,
						},
					}, nil)
					robotStateRepo.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything).Return(nil)
					picCommandSerialRepo.EXPECT().DeletePICSerialCommand(ctx, "123").Return(nil)
				},
				expectedError: false,
			},
			{
				name: "update robot state by command type - drive motor forward",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, robotStateRepo *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{
						Type: model.PICSerialCommandTypeDriveMotor,
						Data: model.PICSerialCommandBatteryDriveMotorData{
							Direction: model.MoveDirectionForward,
							Speed:     50,
							Enable:    true,
						},
					}, nil)
					robotStateRepo.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything).Return(nil)
					picCommandSerialRepo.EXPECT().DeletePICSerialCommand(ctx, "123").Return(nil)
				},
				expectedError: false,
			},
			{
				name: "update robot state by command type - drive motor backward",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, robotStateRepo *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{
						Type: model.PICSerialCommandTypeDriveMotor,
						Data: model.PICSerialCommandBatteryDriveMotorData{
							Direction: model.MoveDirectionBackward,
							Speed:     50,
							Enable:    true,
						},
					}, nil)
					robotStateRepo.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything).Return(nil)
					picCommandSerialRepo.EXPECT().DeletePICSerialCommand(ctx, "123").Return(nil)
				},
				expectedError: false,
			},
			{
				name: "invalid drive motor direction",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, robotStateRepo *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{
						Type: model.PICSerialCommandTypeDriveMotor,
						Data: model.PICSerialCommandBatteryDriveMotorData{
							Direction: 99, // Invalid direction
							Speed:     50,
							Enable:    true,
						},
					}, nil)
					robotStateRepo.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, nil)
				},
				expectedError: true,
			},
			{
				name: "invalid command type",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, robotStateRepo *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{
						Type: 99, // Unknown command type
						Data: nil,
					}, nil)
					robotStateRepo.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, nil)
				},
				expectedError: true,
			},
			{
				name: "invalid command data type battery charge",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, robotStateRepo *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{
						Type: model.PICSerialCommandTypeBatteryCharge,
						Data: model.PICSerialCommandBatteryDischargeData{}, // Invalid data type
					}, nil)
					robotStateRepo.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, nil)
				},
				expectedError: true,
			},
			{
				name: "invalid command data type battery discharge",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, robotStateRepo *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{
						Type: model.PICSerialCommandTypeBatteryDischarge,
						Data: model.PICSerialCommandBatteryChargeData{}, // Invalid data type
					}, nil)
					robotStateRepo.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, nil)
				},
				expectedError: true,
			},
			{
				name: "invalid command data type battery lift motor",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, robotStateRepo *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{
						Type: model.PICSerialCommandTypeLiftMotor,
						Data: model.PICSerialCommandBatteryChargeData{}, // Invalid data type
					}, nil)
					robotStateRepo.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, nil)
				},
				expectedError: true,
			},
			{
				name: "invalid command data type battery drive motor",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, robotStateRepo *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{
						Type: model.PICSerialCommandTypeDriveMotor,
						Data: model.PICSerialCommandBatteryChargeData{}, // Invalid data type
					}, nil)
					robotStateRepo.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, nil)
				},
				expectedError: true,
			},
			{
				name: "update robot state failed",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, robotStateRepo *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{
						Type: model.PICSerialCommandTypeBatteryCharge,
						Data: model.PICSerialCommandBatteryChargeData{
							CurrentLimit: 1,
							Enable:       true,
						},
					}, nil)
					robotStateRepo.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything).Return(assert.AnError)
				},
				expectedError: true,
			},
			{
				name: "delete command failed",
				params: service.ProcessSerialCommandACK{
					ID:      "123",
					Success: true,
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, robotStateRepo *mocks.FakeRobotStateRepository) {
					picCommandSerialRepo.EXPECT().GetPICSerialCommand(ctx, "123").Return(model.PICSerialCommand{
						Type: model.PICSerialCommandTypeBatteryCharge,
						Data: model.PICSerialCommandBatteryChargeData{
							CurrentLimit: 1,
							Enable:       true,
						},
					}, nil)
					robotStateRepo.EXPECT().GetRobotState(ctx).Return(model.RobotState{}, nil)
					robotStateRepo.EXPECT().UpdateRobotState(ctx, mock.Anything).Return(nil)
					picCommandSerialRepo.EXPECT().DeletePICSerialCommand(ctx, "123").Return(assert.AnError)
				},
				expectedError: true,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				picCommandSerialRepo := mocks.NewFakePICSerialCommandRepository(t)
				robotStateRepo := mocks.NewFakeRobotStateRepository(t)
				s := NewPICService(robotStateRepo, picCommandSerialRepo, validator)

				tc.mock(picCommandSerialRepo, robotStateRepo)

				err := s.ProcessSerialCommandACK(ctx, tc.params)
				if tc.expectedError {
					if err == nil {
						t.Errorf("expected error, got nil")
					}
					return
				}
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			})
		}
	})
}
