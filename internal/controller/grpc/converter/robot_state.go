package converter

import (
	raybotv1 "buf.build/gen/go/tbe-team/raybot-api/protocolbuffers/go/raybot/v1"

	"github.com/tbe-team/raybot/internal/model"
)

func ToGetRobotStateResponse(state model.RobotState) *raybotv1.GetRobotStateResponse {
	return &raybotv1.GetRobotStateResponse{
		RobotState: &raybotv1.RobotState{
			BatteryState: &raybotv1.BatteryState{
				Current:          uint32(state.Battery.Current),
				ChargePercentage: uint32(state.Battery.Percent),
				Voltage:          uint32(state.Battery.Voltage),
				Temperature:      uint32(state.Battery.Temp),
			},
			DistanceSensorData: &raybotv1.DistanceSensorData{
				FrontDistance: uint32(state.DistanceSensor.FrontDistance),
				BackDistance:  uint32(state.DistanceSensor.BackDistance),
				DownDistance:  uint32(state.DistanceSensor.DownDistance),
			},
			LiftMotorState: &raybotv1.LiftMotorState{
				CurrentPosition: uint32(state.LiftMotor.CurrentPosition),
				TargetPosition:  uint32(state.LiftMotor.TargetPosition),
				IsRunning:       state.LiftMotor.IsRunning,
				Enabled:         state.LiftMotor.Enabled,
			},
			DriveMotorState: &raybotv1.DriveMotorState{
				Direction: robotStateToDriveMotorDirection(state.DriveMotor.Direction),
				Speed:     uint32(state.DriveMotor.Speed),
				IsRunning: state.DriveMotor.IsRunning,
				Enabled:   state.DriveMotor.Enabled,
			},
		},
	}
}

func robotStateToDriveMotorDirection(direction model.DriveMotorDirection) raybotv1.DriveMotorState_Direction {
	switch direction {
	case model.DriveMotorDirectionForward:
		return raybotv1.DriveMotorState_DIRECTION_FORWARD
	case model.DriveMotorDirectionBackward:
		return raybotv1.DriveMotorState_DIRECTION_BACKWARD
	default:
		return raybotv1.DriveMotorState_DIRECTION_UNSPECIFIED
	}
}
