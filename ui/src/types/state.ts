export interface BatteryState {
  current: number
  temp: number
  voltage: number
  cellVoltages: number[]
  percent: number
  fault: number
  health: number
  updatedAt: Date
}

export interface ChargeState {
  currentLimit: number
  enabled: boolean
  updatedAt: Date
}

export interface DischargeState {
  currentLimit: number
  enabled: boolean
  updatedAt: Date
}

export interface DistanceSensorState {
  frontDistance: number
  backDistance: number
  downDistance: number
  updatedAt: Date
}

export interface LiftMotorState {
  currentPosition: number
  targetPosition: number
  isRunning: boolean
  enabled: boolean
  updatedAt: Date
}

export type DriveMotorDirection = 'Forward' | 'Backward'

export interface DriveMotorState {
  direction: DriveMotorDirection
  speed: number
  isRunning: boolean
  enabled: boolean
  updatedAt: Date
}

export interface RobotState {
  battery: BatteryState
  charge: ChargeState
  discharge: DischargeState
  distanceSensor: DistanceSensorState
  liftMotor: LiftMotorState
  driveMotor: DriveMotorState
}
export interface LocationState {
  targetLocation: string
  currentLocation: string
  updatedAt: Date
}
