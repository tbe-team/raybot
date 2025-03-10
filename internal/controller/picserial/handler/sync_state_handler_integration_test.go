package handler_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tbe-team/raybot/internal/controller/picserial/handler"
	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository/repoimpl"
	"github.com/tbe-team/raybot/internal/service/serviceimpl"
	"github.com/tbe-team/raybot/pkg/validator"
)

func TestIntegrationSyncStateHandler_Handle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	validator := validator.New()

	tests := []struct {
		name     string
		message  handler.SyncStateMessage
		expected func(state model.RobotState) bool
	}{
		{
			name: "battery update",
			message: handler.SyncStateMessage{
				StateType: handler.SyncStateTypeBattery,
				Data:      []byte(`{"current": 500, "temp": 25, "voltage": 12000, "cell_voltages": [4000, 4000, 4000], "percent": 100, "fault": 0, "health": 100}`),
			},
			expected: func(state model.RobotState) bool {
				return state.Battery.Current == 500 &&
					state.Battery.Temp == 25 &&
					state.Battery.Voltage == 12000 &&
					len(state.Battery.CellVoltages) == 3 &&
					state.Battery.Percent == 100 &&
					state.Battery.Fault == 0 &&
					state.Battery.Health == 100
			},
		},
		{
			name: "charge update",
			message: handler.SyncStateMessage{
				StateType: handler.SyncStateTypeCharge,
				Data:      []byte(`{"current_limit": 1000, "enabled": 1}`),
			},
			expected: func(state model.RobotState) bool {
				return state.Charge.CurrentLimit == 1000 &&
					state.Charge.Enabled == true
			},
		},
		{
			name: "discharge update",
			message: handler.SyncStateMessage{
				StateType: handler.SyncStateTypeDischarge,
				Data:      []byte(`{"current_limit": 2000, "enabled": 0}`),
			},
			expected: func(state model.RobotState) bool {
				return state.Discharge.CurrentLimit == 2000 &&
					state.Discharge.Enabled == false
			},
		},
		{
			name: "distance sensor update",
			message: handler.SyncStateMessage{
				StateType: handler.SyncStateTypeDistanceSensor,
				Data:      []byte(`{"front_distance": 100, "back_distance": 200, "down_distance": 50}`),
			},
			expected: func(state model.RobotState) bool {
				return state.DistanceSensor.FrontDistance == 100 &&
					state.DistanceSensor.BackDistance == 200 &&
					state.DistanceSensor.DownDistance == 50
			},
		},
		{
			name: "lift motor update",
			message: handler.SyncStateMessage{
				StateType: handler.SyncStateTypeLiftMotor,
				Data:      []byte(`{"current_position": 100, "target_position": 200, "is_running": 1, "enabled": 1}`),
			},
			expected: func(state model.RobotState) bool {
				return state.LiftMotor.CurrentPosition == 100 &&
					state.LiftMotor.TargetPosition == 200 &&
					state.LiftMotor.IsRunning == true &&
					state.LiftMotor.Enabled == true
			},
		},
		{
			name: "drive motor update",
			message: handler.SyncStateMessage{
				StateType: handler.SyncStateTypeDriveMotor,
				Data:      []byte(`{"direction": 1, "speed": 50, "is_running": 1, "enabled": 1}`),
			},
			expected: func(state model.RobotState) bool {
				return state.DriveMotor.Direction == model.DriveMotorDirectionBackward &&
					state.DriveMotor.Speed == 50 &&
					state.DriveMotor.IsRunning == true &&
					state.DriveMotor.Enabled == true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			robotStateRepo := repoimpl.NewRobotStateRepository()
			robotService := serviceimpl.NewRobotService(robotStateRepo, validator)
			handler := handler.NewSyncStateHandler(robotService)

			handler.Handle(ctx, tt.message)

			state, err := robotStateRepo.GetRobotState(ctx)
			if err != nil {
				t.Fatalf("failed to get robot state: %v", err)
			}
			assert.True(t, tt.expected(state))
		})
	}
}
