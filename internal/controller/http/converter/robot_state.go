package converter

import (
	"github.com/tbe-team/raybot/internal/controller/http/oas/gen"
	"github.com/tbe-team/raybot/internal/model"
)

func ConvertRobotStateToResponse(state model.RobotState) gen.RobotStateResponse {
	return gen.RobotStateResponse{
		Battery: gen.BatteryState{
			Current:      state.Battery.Current,
			Temp:         state.Battery.Temp,
			Voltage:      state.Battery.Voltage,
			CellVoltages: state.Battery.CellVoltages,
			Percent:      state.Battery.Percent,
			Fault:        state.Battery.Fault,
			Health:       state.Battery.Health,
			UpdatedAt:    state.Battery.UpdatedAt,
		},
		Charge: gen.ChargeState{
			CurrentLimit: state.Charge.CurrentLimit,
			Enabled:      state.Charge.Enabled,
			UpdatedAt:    state.Charge.UpdatedAt,
		},
		Discharge: gen.DischargeState{
			CurrentLimit: state.Discharge.CurrentLimit,
			Enabled:      state.Discharge.Enabled,
			UpdatedAt:    state.Discharge.UpdatedAt,
		},
		DistanceSensor: gen.DistanceSensorState{
			FrontDistance: state.DistanceSensor.FrontDistance,
			BackDistance:  state.DistanceSensor.BackDistance,
			DownDistance:  state.DistanceSensor.DownDistance,
			UpdatedAt:     state.DistanceSensor.UpdatedAt,
		},
		LiftMotor: gen.LiftMotorState{
			CurrentPosition: state.LiftMotor.CurrentPosition,
			TargetPosition:  state.LiftMotor.TargetPosition,
			IsRunning:       state.LiftMotor.IsRunning,
			Enabled:         state.LiftMotor.Enabled,
			UpdatedAt:       state.LiftMotor.UpdatedAt,
		},
		DriveMotor: gen.DriveMotorState{
			Direction: state.DriveMotor.Direction.String(),
			Speed:     state.DriveMotor.Speed,
			IsRunning: state.DriveMotor.IsRunning,
			Enabled:   state.DriveMotor.Enabled,
			UpdatedAt: state.DriveMotor.UpdatedAt,
		},
		Location: gen.LocationState{
			CurrentLocation: state.Location.CurrentLocation,
			UpdatedAt:       state.Location.UpdatedAt,
		},
		Cargo: gen.CargoState{
			IsOpen:         state.Cargo.IsOpen,
			QrCode:         state.Cargo.QRCode,
			BottomDistance: state.Cargo.BottomDistance,
			UpdatedAt:      state.Cargo.UpdatedAt,
		},
		CargoDoorMotor: gen.CargoDoorMotorState{
			Direction: state.CargoDoorMotor.Direction.String(),
			Speed:     state.CargoDoorMotor.Speed,
			IsRunning: state.CargoDoorMotor.IsRunning,
			Enabled:   state.CargoDoorMotor.Enabled,
			UpdatedAt: state.CargoDoorMotor.UpdatedAt,
		},
	}
}
