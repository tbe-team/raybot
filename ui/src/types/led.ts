const _mode = ['OFF', 'ON', 'BLINK'] as const
type LedMode = (typeof _mode)[number]

interface LedConnection {
  connected: boolean
  lastConnectedAt?: string
  error?: string
}

interface LedState {
  mode: LedMode
  updatedAt: string
}

export interface Led {
  connection: LedConnection
  state: LedState
}
