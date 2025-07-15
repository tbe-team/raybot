export interface SystemInfo {
  localIp: string
  cpuUsage: number
  memoryUsage: number
  totalMemory: number
  uptime: number
}

export interface SystemStatusResponse {
  status: SystemStatusType
}

export type SystemStatusType = 'NORMAL' | 'ERROR'
