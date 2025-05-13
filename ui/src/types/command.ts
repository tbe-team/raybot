export interface StopMovementInputs {}
export interface MoveForwardInputs {}
export interface MoveBackwardInputs {}
export interface MoveToInputs {
  location: string
  direction: 'FORWARD' | 'BACKWARD'
}
export interface CargoOpenInputs {}
export interface CargoCloseInputs {}
export interface CargoLiftInputs {}
export interface CargoLowerInputs {}
export interface CargoCheckQRInputs {
  qrCode: string
}
export interface ScanLocationInputs {}
export interface WaitInputs {
  durationMs: number
}
export interface StopMovementOutputs {}
export interface MoveForwardOutputs {}
export interface MoveBackwardOutputs {}
export interface MoveToOutputs {}
export interface CargoOpenOutputs {}
export interface CargoCloseOutputs {}
export interface CargoLiftOutputs {}
export interface CargoLowerOutputs {}
export interface CargoCheckQROutputs {}
export interface ScanLocationOutputs {
  locations: Location[]
}
export interface Location {
  location: string
  scannedAt: string
}
export interface WaitOutputs {}
export interface CommandInputMap {
  STOP_MOVEMENT: StopMovementInputs
  MOVE_FORWARD: MoveForwardInputs
  MOVE_BACKWARD: MoveBackwardInputs
  MOVE_TO: MoveToInputs
  CARGO_OPEN: CargoOpenInputs
  CARGO_CLOSE: CargoCloseInputs
  CARGO_LIFT: CargoLiftInputs
  CARGO_LOWER: CargoLowerInputs
  CARGO_CHECK_QR: CargoCheckQRInputs
  SCAN_LOCATION: ScanLocationInputs
  WAIT: WaitInputs
}
export interface CommandOutputMap {
  STOP_MOVEMENT: StopMovementOutputs
  MOVE_FORWARD: MoveForwardOutputs
  MOVE_BACKWARD: MoveBackwardOutputs
  MOVE_TO: MoveToOutputs
  CARGO_OPEN: CargoOpenOutputs
  CARGO_CLOSE: CargoCloseOutputs
  CARGO_LIFT: CargoLiftOutputs
  CARGO_LOWER: CargoLowerOutputs
  CARGO_CHECK_QR: CargoCheckQROutputs
  SCAN_LOCATION: ScanLocationOutputs
  WAIT: WaitOutputs
}

export const CommandTypeValues = [
  'STOP_MOVEMENT',
  'MOVE_FORWARD',
  'MOVE_BACKWARD',
  'MOVE_TO',
  'CARGO_OPEN',
  'CARGO_CLOSE',
  'CARGO_LIFT',
  'CARGO_LOWER',
  'CARGO_CHECK_QR',
  'SCAN_LOCATION',
  'WAIT',
] as const
export type CommandType = typeof CommandTypeValues[number]

export type CommandStatus = 'QUEUED' | 'PROCESSING' | 'SUCCEEDED' | 'FAILED' | 'CANCELED'

export type CommandSource = 'CLOUD' | 'APP'

export interface Command<T extends CommandType = CommandType> {
  id: number
  type: T
  status: CommandStatus
  source: CommandSource
  inputs: CommandInputMap[T]
  outputs: CommandOutputMap[T]
  error?: string
  completedAt?: string
  startedAt?: string
  createdAt: string
  updatedAt: string
}
