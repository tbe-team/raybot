import type { Cargo, CargoDoorMotor } from './cargo'

export interface BatteryState {
  current: number
  temp: number
  voltage: number
  cellVoltages: number[]
  percent: number
  fault: number
  health: number
  updatedAt: string
}

export interface ChargeState {
  currentLimit: number
  enabled: boolean
  updatedAt: string
}

export interface DischargeState {
  currentLimit: number
  enabled: boolean
  updatedAt: string
}

export interface DistanceSensorState {
  frontDistance: number
  backDistance: number
  downDistance: number
  updatedAt: string
}

export interface LiftMotorState {
  currentPosition: number
  targetPosition: number
  isRunning: boolean
  enabled: boolean
  updatedAt: string
}

export interface DriveMotorState {
  direction: 'FORWARD' | 'BACKWARD'
  speed: number
  isRunning: boolean
  enabled: boolean
  updatedAt: string
}

export interface LocationState {
  currentLocation: string
  updatedAt: string
}

export interface RobotState {
  battery: BatteryState
  charge: ChargeState
  discharge: DischargeState
  distanceSensor: DistanceSensorState
  liftMotor: LiftMotorState
  driveMotor: DriveMotorState
  location: LocationState
  cargo: Cargo
  cargoDoorMotor: CargoDoorMotor
}
