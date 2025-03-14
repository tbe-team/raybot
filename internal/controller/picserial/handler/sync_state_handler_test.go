package handler_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tbe-team/raybot/internal/controller/picserial/handler"
	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/service/mocks"
	"github.com/tbe-team/raybot/pkg/log"
)

func TestUnmarshalSyncStateType(t *testing.T) {
	tests := []struct {
		name      string
		msg       []byte
		shouldErr bool
		want      handler.SyncStateType
	}{
		{
			name:      "sync state type battery",
			msg:       []byte(`{"state_type": 0}`),
			shouldErr: false,
			want:      handler.SyncStateTypeBattery,
		},
		{
			name:      "sync state type charge",
			msg:       []byte(`{"state_type": 1}`),
			shouldErr: false,
			want:      handler.SyncStateTypeCharge,
		},
		{
			name:      "sync state type discharge",
			msg:       []byte(`{"state_type": 2}`),
			shouldErr: false,
			want:      handler.SyncStateTypeDischarge,
		},
		{
			name:      "sync state type distance sensor",
			msg:       []byte(`{"state_type": 3}`),
			shouldErr: false,
			want:      handler.SyncStateTypeDistanceSensor,
		},
		{
			name:      "sync state type lift motor",
			msg:       []byte(`{"state_type": 4}`),
			shouldErr: false,
			want:      handler.SyncStateTypeLiftMotor,
		},
		{
			name:      "sync state type drive motor",
			msg:       []byte(`{"state_type": 5}`),
			shouldErr: false,
			want:      handler.SyncStateTypeDriveMotor,
		},
		{
			name:      "invalid sync state type",
			msg:       []byte(`{"state_type": 6}`),
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var temp struct {
				StateType handler.SyncStateType `json:"state_type"`
			}
			err := json.Unmarshal(tt.msg, &temp)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if temp.StateType != tt.want {
					t.Errorf("expected %v, got %v", tt.want, temp.StateType)
				}
			}
		})
	}
}

func TestSyncStateHandler_Handle(t *testing.T) {
	log := log.NewNopLogger()
	ctx := context.Background()

	tests := []struct {
		name          string
		msg           handler.SyncStateMessage
		mockSetup     func(*mocks.FakeRobotService)
		expectedError bool
	}{
		{
			name: "process battery state",
			msg: handler.SyncStateMessage{
				StateType: handler.SyncStateTypeBattery,
				Data:      []byte(`{"current": 500, "temp": 25, "voltage": 12000, "cell_voltages": [4000, 4000, 4000], "percent": 100, "fault": 0, "health": 100}`),
			},
			mockSetup: func(robotService *mocks.FakeRobotService) {
				robotService.EXPECT().UpdateRobotState(ctx, mock.MatchedBy(func(params service.UpdateRobotStateParams) bool {
					return params.SetBattery &&
						params.Battery.Current == 500 &&
						params.Battery.Temp == 25 &&
						params.Battery.Voltage == 12000 &&
						len(params.Battery.CellVoltages) == 3 &&
						params.Battery.Percent == 100 &&
						params.Battery.Fault == 0 &&
						params.Battery.Health == 100
				})).Return(model.RobotState{}, nil)
			},
			expectedError: false,
		},
		{
			name: "process charge state",
			msg: handler.SyncStateMessage{
				StateType: handler.SyncStateTypeCharge,
				Data:      []byte(`{"current_limit": 1000, "enabled": 1}`),
			},
			mockSetup: func(robotService *mocks.FakeRobotService) {
				robotService.EXPECT().UpdateRobotState(ctx, mock.MatchedBy(func(params service.UpdateRobotStateParams) bool {
					return params.SetCharge &&
						params.Charge.CurrentLimit == 1000 &&
						params.Charge.Enabled == true
				})).Return(model.RobotState{}, nil)
			},
			expectedError: false,
		},
		{
			name: "process discharge state",
			msg: handler.SyncStateMessage{
				StateType: handler.SyncStateTypeDischarge,
				Data:      []byte(`{"current_limit": 2000, "enabled": 0}`),
			},
			mockSetup: func(robotService *mocks.FakeRobotService) {
				robotService.EXPECT().UpdateRobotState(ctx, mock.MatchedBy(func(params service.UpdateRobotStateParams) bool {
					return params.SetDischarge &&
						params.Discharge.CurrentLimit == 2000 &&
						params.Discharge.Enabled == false
				})).Return(model.RobotState{}, nil)
			},
			expectedError: false,
		},
		{
			name: "process distance sensor state",
			msg: handler.SyncStateMessage{
				StateType: handler.SyncStateTypeDistanceSensor,
				Data:      []byte(`{"front_distance": 100, "back_distance": 200, "down_distance": 50}`),
			},
			mockSetup: func(robotService *mocks.FakeRobotService) {
				robotService.EXPECT().UpdateRobotState(ctx, mock.MatchedBy(func(params service.UpdateRobotStateParams) bool {
					return params.SetDistanceSensor &&
						params.DistanceSensor.FrontDistance == 100 &&
						params.DistanceSensor.BackDistance == 200 &&
						params.DistanceSensor.DownDistance == 50
				})).Return(model.RobotState{}, nil)
			},
			expectedError: false,
		},
		{
			name: "process lift motor state",
			msg: handler.SyncStateMessage{
				StateType: handler.SyncStateTypeLiftMotor,
				Data:      []byte(`{"current_position": 100, "target_position": 200, "is_running": 1, "enabled": 1}`),
			},
			mockSetup: func(robotService *mocks.FakeRobotService) {
				robotService.EXPECT().UpdateRobotState(ctx, mock.MatchedBy(func(params service.UpdateRobotStateParams) bool {
					return params.SetLiftMotor &&
						params.LiftMotor.CurrentPosition == 100 &&
						params.LiftMotor.TargetPosition == 200 &&
						params.LiftMotor.IsRunning == true &&
						params.LiftMotor.Enabled == true
				})).Return(model.RobotState{}, nil)
			},
			expectedError: false,
		},
		{
			name: "process drive motor state",
			msg: handler.SyncStateMessage{
				StateType: handler.SyncStateTypeDriveMotor,
				Data:      []byte(`{"direction": 1, "speed": 50, "is_running": 1, "enabled": 1}`),
			},
			mockSetup: func(robotService *mocks.FakeRobotService) {
				robotService.EXPECT().UpdateRobotState(ctx, mock.MatchedBy(func(params service.UpdateRobotStateParams) bool {
					return params.SetDriveMotor &&
						params.DriveMotor.Direction == model.DriveMotorDirectionBackward &&
						params.DriveMotor.Speed == 50 &&
						params.DriveMotor.IsRunning == true &&
						params.DriveMotor.Enabled == true
				})).Return(model.RobotState{}, nil)
			},
			expectedError: false,
		},
		{
			name: "invalid battery data",
			msg: handler.SyncStateMessage{
				StateType: handler.SyncStateTypeBattery,
				Data:      []byte(`{"invalid": "json"}`),
			},
			mockSetup: func(robotService *mocks.FakeRobotService) {
				robotService.EXPECT().UpdateRobotState(ctx, mock.MatchedBy(func(params service.UpdateRobotStateParams) bool {
					return params.SetBattery
				})).Return(model.RobotState{}, nil)
			},
			expectedError: false,
		},
		{
			name: "update robot state error",
			msg: handler.SyncStateMessage{
				StateType: handler.SyncStateTypeBattery,
				Data:      []byte(`{"current": 500, "temp": 25, "voltage": 12000, "cell_voltages": [4000, 4000, 4000], "percent": 100, "fault": 0, "health": 100}`),
			},
			mockSetup: func(robotService *mocks.FakeRobotService) {
				robotService.EXPECT().UpdateRobotState(ctx, mock.Anything).Return(model.RobotState{}, assert.AnError)
			},
			expectedError: true,
		},
		{
			name: "unknown state type",
			msg: handler.SyncStateMessage{
				StateType: 99,
				Data:      []byte(`{}`),
			},
			mockSetup: func(_ *mocks.FakeRobotService) {
			},
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			robotService := mocks.NewFakeRobotService(t)
			syncStateHandler := handler.NewSyncStateHandler(robotService, log)

			tc.mockSetup(robotService)

			syncStateHandler.Handle(ctx, tc.msg)
		})
	}
}
