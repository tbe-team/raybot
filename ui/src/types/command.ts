export const CommandTypeValues = ['MOVE_TO_LOCATION', 'LIFT_CARGO', 'DROP_CARGO', 'OPEN_CARGO', 'CLOSE_CARGO'] as const
export type CommandType = typeof CommandTypeValues[number]

export type CommandStatus = 'IN_PROGRESS' | 'SUCCEEDED' | 'FAILED'

export type CommandSource = 'MANUAL' | 'CLOUD'

export interface Command {
  id: string
  type: CommandType
  status: CommandStatus
  source: CommandSource
  inputs: Record<string, unknown>
  error?: string | null
  createdAt: string
  completedAt?: string | null
}
