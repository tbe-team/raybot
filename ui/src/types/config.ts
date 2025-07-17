export type LogLevel = 'DEBUG' | 'INFO' | 'WARN' | 'ERROR'
export type LogFormat = 'JSON' | 'TEXT'
export interface SerialPort {
  items: { port: string }[]
}

export interface LogConsoleConfig {
  enable: boolean
  level: LogLevel
  format: LogFormat
}

export interface LogFileConfig {
  enable: boolean
  path: string
  rotationCount: number
  format: LogFormat
  level: LogLevel
}

export interface LogConfig {
  file: LogFileConfig
  console: LogConsoleConfig
}

export interface LEDConfig {
  pin: string
}

export interface LEDsConfig {
  system: LEDConfig
  alert: LEDConfig
}

export interface HardwareConfig {
  esp: ESPConfig
  pic: PICConfig
  leds: LEDsConfig
}

export interface ESPConfig {
  serial: SerialConfig
  enableAck: boolean
  commandAckTimeout: number
}

export interface PICConfig {
  serial: SerialConfig
  enableAck: boolean
  commandAckTimeout: number
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
  enable: boolean
  address: string
  token: string
}

export interface HTTPConfig {
  port: number
  swagger: boolean
}

export interface WifiConfig {
  ap: APConfig
  sta: STAConfig
}

export interface APConfig {
  enable: boolean
  ssid: string
  password: string
  ip: string
}

export interface STAConfig {
  enable: boolean
  ssid: string
  password: string
  ip: string
}

export interface BatteryMonitoringConfig {
  voltageLow: {
    enable: boolean
    threshold: number
  }
  voltageHigh: {
    enable: boolean
    threshold: number
  }
  cellVoltageHigh: {
    enable: boolean
    threshold: number
  }
  cellVoltageLow: {
    enable: boolean
    threshold: number
  }
  cellVoltageDiff: {
    enable: boolean
    threshold: number
  }
  currentHigh: {
    enable: boolean
    threshold: number
  }
  tempHigh: {
    enable: boolean
    threshold: number
  }
  percentLow: {
    enable: boolean
    threshold: number
  }
  healthLow: {
    enable: boolean
    threshold: number
  }
}
