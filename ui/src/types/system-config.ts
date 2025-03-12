type LogLevel = 'debug' | 'info' | 'warn' | 'error'
type SerialParity = 'none' | 'even' | 'odd'
type LogFormat = 'json' | 'text'
type SerialDataBits = 5 | 6 | 7 | 8
type SerialStopBits = 1 | 1.5 | 2

interface GrpcConfig {
  port: number
}

interface HttpConfig {
  enableSwagger: boolean
}

interface LogConfig {
  level: LogLevel
  format: LogFormat
  addSource: boolean
}

interface SerialConfig {
  port: string
  baudRate: number
  dataBits: SerialDataBits
  stopBits: SerialStopBits
  parity: SerialParity
  readTimeout: number
}

interface PicConfig {
  serial: SerialConfig
}

export interface SystemConfig {
  grpc: GrpcConfig
  http: HttpConfig
  log: LogConfig
  pic: PicConfig
}
