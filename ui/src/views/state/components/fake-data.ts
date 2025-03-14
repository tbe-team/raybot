import type { BatteryState, ChargeState, DischargeState, DistanceSensorState, DriveMotorState, LiftMotorState, LocationState, RobotState } from '@/types/robot-state'

const batteryState: BatteryState = {
  current: 100,
  temp: 25,
  voltage: 120,
  cellVoltages: [12, 12, 12, 12],
  percent: 50,
  fault: 0,
  health: 100,
  updatedAt: '2021-01-01T00:00:00Z',
}

const chargeState: ChargeState = {
  currentLimit: 100,
  enabled: true,
  updatedAt: '2021-01-01T00:00:00Z',
}

const dischargeState: DischargeState = {
  currentLimit: 100,
  enabled: true,
  updatedAt: '2021-01-01T00:00:00Z',
}

const distanceSensorState: DistanceSensorState = {
  frontDistance: 100,
  backDistance: 100,
  downDistance: 100,
  updatedAt: '2021-01-01T00:00:00Z',
}

const liftMotorState: LiftMotorState = {
  currentPosition: 100,
  targetPosition: 100,
  isRunning: true,
  enabled: true,
  updatedAt: '2021-01-01T00:00:00Z',
}

const driveMotorState: DriveMotorState = {
  direction: 'FORWARD',
  speed: 100,
  isRunning: true,
  enabled: true,
  updatedAt: '2021-01-01T00:00:00Z',
}

const locationState: LocationState = {
  currentLocation: 'ABCxyz',
  updatedAt: '2021-01-01T00:00:00Z',
}

export const robotState: RobotState = {
  battery: batteryState,
  charge: chargeState,
  discharge: dischargeState,
  distanceSensor: distanceSensorState,
  liftMotor: liftMotorState,
  driveMotor: driveMotorState,
  location: locationState,
}
