export interface Cargo {
  isOpen: boolean
  qrCode: string
  bottomDistance: number
  hasItem: boolean
  updatedAt: string
}

export interface CargoDoorMotorState {
  direction: 'CLOSE' | 'OPEN'
  speed: number
  isRunning: boolean
  enabled: boolean
  updatedAt: string
}
