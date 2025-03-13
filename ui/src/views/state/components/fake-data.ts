import type { BatteryState, ChargeState, DischargeState, DistanceSensorState, DriveMotorState, LiftMotorState, LocationState, RobotState } from '@/types/state'

const batteryState: BatteryState = {
  current: 100,
  temp: 30,
  voltage: 24.5,
  cellVoltages: [4.1, 4.0, 4.1, 4.2],
  percent: 59,
  fault: 0,
  health: 95,
  updatedAt: new Date(),
}

const chargeState: ChargeState = {
  currentLimit: 10,
  enabled: true,
  updatedAt: new Date(),
}

const dischargeState: DischargeState = {
  currentLimit: 15,
  enabled: false,
  updatedAt: new Date(),
}

const distanceSensorState: DistanceSensorState = {
  frontDistance: 50,
  backDistance: 30,
  downDistance: 10,
  updatedAt: new Date(),
}

const liftMotorState: LiftMotorState = {
  currentPosition: 20,
  targetPosition: 50,
  isRunning: false,
  enabled: true,
  updatedAt: new Date(),
}

const driveMotorState: DriveMotorState = {
  direction: 'Forward',
  speed: 5.5,
  isRunning: true,
  enabled: true,
  updatedAt: new Date(),
}

export const locationState: LocationState = {
  targetLocation: 'TB10',
  currentLocation: 'TB03',
  updatedAt: new Date(),
}

export const robotState: RobotState = {
  battery: batteryState,
  charge: chargeState,
  discharge: dischargeState,
  distanceSensor: distanceSensorState,
  liftMotor: liftMotorState,
  driveMotor: driveMotorState,
}
