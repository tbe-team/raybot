package serviceimpl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	serialmocks "github.com/tbe-team/raybot/internal/controller/picserial/serial/mocks"
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
			params        service.ProcessSerialCommandACKParams
			mock          func(_ *mocks.FakePICSerialCommandRepository, _ *mocks.FakeRobotStateRepository)
			expectedError bool
		}{
			{
				name: "validation failed",
				params: service.ProcessSerialCommandACKParams{
					ID:      "",
					Success: true,
				},
				mock: func(_ *mocks.FakePICSerialCommandRepository, _ *mocks.FakeRobotStateRepository) {
				},
				expectedError: true,
			},
			{
				name: "command execution failed",
				params: service.ProcessSerialCommandACKParams{
					ID:      "123",
					Success: false,
				},
				mock: func(_ *mocks.FakePICSerialCommandRepository, _ *mocks.FakeRobotStateRepository) {
				},
				expectedError: true,
			},
			{
				name: "get pic serial command failed",
				params: service.ProcessSerialCommandACKParams{
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
				params: service.ProcessSerialCommandACKParams{
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
				params: service.ProcessSerialCommandACKParams{
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
				params: service.ProcessSerialCommandACKParams{
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
				params: service.ProcessSerialCommandACKParams{
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
				params: service.ProcessSerialCommandACKParams{
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
				params: service.ProcessSerialCommandACKParams{
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
				params: service.ProcessSerialCommandACKParams{
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
				params: service.ProcessSerialCommandACKParams{
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
				params: service.ProcessSerialCommandACKParams{
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
				params: service.ProcessSerialCommandACKParams{
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
				params: service.ProcessSerialCommandACKParams{
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
				params: service.ProcessSerialCommandACKParams{
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
				params: service.ProcessSerialCommandACKParams{
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
				params: service.ProcessSerialCommandACKParams{
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
				s := NewPICService(robotStateRepo, picCommandSerialRepo, nil, validator)

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

	t.Run("test create serial command", func(t *testing.T) {
		tests := []struct {
			name          string
			params        service.CreateSerialCommandParams
			mock          func(_ *mocks.FakePICSerialCommandRepository, _ *serialmocks.FakeClient)
			expectedError bool
		}{
			{
				name: "validation failed - empty data",
				params: service.CreateSerialCommandParams{
					Data: nil,
				},
				mock: func(_ *mocks.FakePICSerialCommandRepository, _ *serialmocks.FakeClient) {
					// No mocks needed - validation should fail
				},
				expectedError: true,
			},
			{
				name: "create command repository error",
				params: service.CreateSerialCommandParams{
					Data: model.PICSerialCommandBatteryChargeData{
						CurrentLimit: 100,
						Enable:       true,
					},
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, _ *serialmocks.FakeClient) {
					picCommandSerialRepo.EXPECT().CreatePICSerialCommand(ctx, mock.Anything).Return(assert.AnError)
				},
				expectedError: true,
			},
			{
				name: "serial client write error",
				params: service.CreateSerialCommandParams{
					Data: model.PICSerialCommandBatteryChargeData{
						CurrentLimit: 100,
						Enable:       true,
					},
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, picSerialClient *serialmocks.FakeClient) {
					picCommandSerialRepo.EXPECT().CreatePICSerialCommand(ctx, mock.Anything).Return(nil)
					picSerialClient.EXPECT().Write(mock.Anything).Return(assert.AnError)
				},
				expectedError: true,
			},
			{
				name: "battery charge command success",
				params: service.CreateSerialCommandParams{
					Data: model.PICSerialCommandBatteryChargeData{
						CurrentLimit: 100,
						Enable:       true,
					},
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, picSerialClient *serialmocks.FakeClient) {
					picCommandSerialRepo.EXPECT().CreatePICSerialCommand(ctx, mock.MatchedBy(func(cmd model.PICSerialCommand) bool {
						data, ok := cmd.Data.(model.PICSerialCommandBatteryChargeData)
						return ok &&
							cmd.Type == model.PICSerialCommandTypeBatteryCharge &&
							data.CurrentLimit == 100 &&
							data.Enable == true
					})).Return(nil)
					picSerialClient.EXPECT().Write(mock.Anything).Return(nil)
				},
				expectedError: false,
			},
			{
				name: "battery discharge command success",
				params: service.CreateSerialCommandParams{
					Data: model.PICSerialCommandBatteryDischargeData{
						CurrentLimit: 200,
						Enable:       true,
					},
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, picSerialClient *serialmocks.FakeClient) {
					picCommandSerialRepo.EXPECT().CreatePICSerialCommand(ctx, mock.MatchedBy(func(cmd model.PICSerialCommand) bool {
						data, ok := cmd.Data.(model.PICSerialCommandBatteryDischargeData)
						return ok &&
							cmd.Type == model.PICSerialCommandTypeBatteryDischarge &&
							data.CurrentLimit == 200 &&
							data.Enable == true
					})).Return(nil)
					picSerialClient.EXPECT().Write(mock.Anything).Return(nil)
				},
				expectedError: false,
			},
			{
				name: "lift motor command success",
				params: service.CreateSerialCommandParams{
					Data: model.PICSerialCommandBatteryLiftMotorData{
						TargetPosition: 150,
						Enable:         true,
					},
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, picSerialClient *serialmocks.FakeClient) {
					picCommandSerialRepo.EXPECT().CreatePICSerialCommand(ctx, mock.MatchedBy(func(cmd model.PICSerialCommand) bool {
						data, ok := cmd.Data.(model.PICSerialCommandBatteryLiftMotorData)
						return ok &&
							cmd.Type == model.PICSerialCommandTypeLiftMotor &&
							data.TargetPosition == 150 &&
							data.Enable == true
					})).Return(nil)
					picSerialClient.EXPECT().Write(mock.Anything).Return(nil)
				},
				expectedError: false,
			},
			{
				name: "drive motor forward command success",
				params: service.CreateSerialCommandParams{
					Data: model.PICSerialCommandBatteryDriveMotorData{
						Direction: model.MoveDirectionForward,
						Speed:     75,
						Enable:    true,
					},
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, picSerialClient *serialmocks.FakeClient) {
					picCommandSerialRepo.EXPECT().CreatePICSerialCommand(ctx, mock.MatchedBy(func(cmd model.PICSerialCommand) bool {
						data, ok := cmd.Data.(model.PICSerialCommandBatteryDriveMotorData)
						return ok &&
							cmd.Type == model.PICSerialCommandTypeDriveMotor &&
							data.Direction == model.MoveDirectionForward &&
							data.Speed == 75 &&
							data.Enable == true
					})).Return(nil)
					picSerialClient.EXPECT().Write(mock.Anything).Return(nil)
				},
				expectedError: false,
			},
			{
				name: "drive motor backward command success",
				params: service.CreateSerialCommandParams{
					Data: model.PICSerialCommandBatteryDriveMotorData{
						Direction: model.MoveDirectionBackward,
						Speed:     50,
						Enable:    true,
					},
				},
				mock: func(picCommandSerialRepo *mocks.FakePICSerialCommandRepository, picSerialClient *serialmocks.FakeClient) {
					picCommandSerialRepo.EXPECT().CreatePICSerialCommand(ctx, mock.MatchedBy(func(cmd model.PICSerialCommand) bool {
						data, ok := cmd.Data.(model.PICSerialCommandBatteryDriveMotorData)
						return ok &&
							cmd.Type == model.PICSerialCommandTypeDriveMotor &&
							data.Direction == model.MoveDirectionBackward &&
							data.Speed == 50 &&
							data.Enable == true
					})).Return(nil)
					picSerialClient.EXPECT().Write(mock.Anything).Return(nil)
				},
				expectedError: false,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				picCommandSerialRepo := mocks.NewFakePICSerialCommandRepository(t)
				picSerialClient := serialmocks.NewFakeClient(t)
				s := NewPICService(mocks.NewFakeRobotStateRepository(t), picCommandSerialRepo, picSerialClient, validator)

				tc.mock(picCommandSerialRepo, picSerialClient)

				err := s.CreateSerialCommand(ctx, tc.params)
				if tc.expectedError {
					assert.Error(t, err)
					return
				}
				assert.NoError(t, err)
			})
		}
	})
}
