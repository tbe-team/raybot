export interface Cargo {
  isOpen: boolean
  qrCode: string
  bottomDistance: number
  updatedAt: string
}

export interface CargoDoorMotor {
  direction: 'CLOSE' | 'OPEN'
  speed: number
  isRunning: boolean
  enabled: boolean
  updatedAt: string
}
