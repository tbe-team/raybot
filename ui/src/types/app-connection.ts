export interface CloudConnection {
  connected: boolean
  lastConnectedAt?: string
  uptime: number
  error?: string
}

export interface ESPSerialConnection {
  connected: boolean
  lastConnectedAt?: string
  error?: string
}

export interface PICSerialConnection {
  connected: boolean
  lastConnectedAt?: string
  error?: string
}

export interface RFIDUSBConnection {
  connected: boolean
  lastConnectedAt?: string
  error?: string
}

export interface AppConnection {
  cloudConnection: CloudConnection
  espSerialConnection: ESPSerialConnection
  picSerialConnection: PICSerialConnection
  rfidUsbConnection: RFIDUSBConnection
}
