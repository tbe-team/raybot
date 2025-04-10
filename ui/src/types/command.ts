export interface StopInputs {}
export interface MoveForwardInputs {}
export interface MoveBackwardInputs {}
export interface MoveToInputs {
  location: string
}
export interface CargoOpenInputs {}
export interface CargoCloseInputs {}
export interface CargoLiftInputs {}
export interface CargoLowerInputs {}
export interface CargoCheckQRInputs {
  qrCode: string
}

export interface CommandInputMap {
  STOP: StopInputs
  MOVE_FORWARD: MoveForwardInputs
  MOVE_BACKWARD: MoveBackwardInputs
  MOVE_TO: MoveToInputs
  CARGO_OPEN: CargoOpenInputs
  CARGO_CLOSE: CargoCloseInputs
  CARGO_LIFT: CargoLiftInputs
  CARGO_LOWER: CargoLowerInputs
  CARGO_CHECK_QR: CargoCheckQRInputs
}

export const CommandTypeValues = [
  'STOP',
  'MOVE_FORWARD',
  'MOVE_BACKWARD',
  'MOVE_TO',
  'CARGO_OPEN',
  'CARGO_CLOSE',
  'CARGO_LIFT',
  'CARGO_LOWER',
  'CARGO_CHECK_QR',
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
  error?: string | null
  completedAt?: string | null
  createdAt: string
  updatedAt: string
}
