export type LogLevel = 'DEBUG' | 'INFO' | 'WARN' | 'ERROR'
export type LogFormat = 'JSON' | 'TEXT'
export interface SerialPort {
  items: { port: string }[]
}

export interface LogConfig {
  level: LogLevel
  format: LogFormat
  addSource: boolean
}

export interface HardwareConfig {
  esp: ESPConfig
  pic: PICConfig
}

export interface ESPConfig {
  serial: SerialConfig
}

export interface PICConfig {
  serial: SerialConfig
}

export type Parity = 'NONE' | 'EVEN' | 'ODD'
export type DataBits = 5 | 6 | 7 | 8
export type StopBits = 1 | 1.5 | 2

export interface SerialConfig {
  port: string
  baudRate: number
  parity: Parity
  dataBits: DataBits
  stopBits: StopBits
  readTimeout: number
}

export interface CloudConfig {
  address: string
  token: string
}

export interface GRPCConfig {
  port: number
  enable: boolean
}

export interface HTTPConfig {
  port: number
  swagger: boolean
}
