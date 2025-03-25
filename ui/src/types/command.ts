export type CommandType = 'MOVE_TO_LOCATION'

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
