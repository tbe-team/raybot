import type { CommandConfig } from '@/types/command-config'
import type { BatteryMonitoringConfig, CloudConfig, HardwareConfig, HTTPConfig, LogConfig, WifiConfig } from '@/types/config'
import http from '@/lib/http'

const configAPI = {
  getLogConfig: (): Promise<LogConfig> => http.get('/configs/log'),
  updateLogConfig: (config: LogConfig): Promise<void> => http.put('/configs/log', config),

  getHardwareConfig: (): Promise<HardwareConfig> => http.get('/configs/hardware'),
  updateHardwareConfig: (config: HardwareConfig): Promise<void> => http.put('/configs/hardware', config),

  getCloudConfig: (): Promise<CloudConfig> => http.get('/configs/cloud'),
  updateCloudConfig: (config: CloudConfig): Promise<void> => http.put('/configs/cloud', config),

  getHttpConfig: (): Promise<HTTPConfig> => http.get('/configs/http'),
  updateHttpConfig: (config: HTTPConfig): Promise<void> => http.put('/configs/http', config),

  getWifiConfig: (): Promise<WifiConfig> => http.get('/configs/wifi'),
  updateWifiConfig: (config: WifiConfig): Promise<void> => http.put('/configs/wifi', config),

  getCommandConfig: (): Promise<CommandConfig> => http.get('/configs/command'),
  updateCommandConfig: (config: CommandConfig): Promise<void> => http.put('/configs/command', config),

  getBatteryMonitoringConfig: (): Promise<BatteryMonitoringConfig> => http.get('/configs/monitoring/battery'),
  updateBatteryMonitoringConfig: (config: BatteryMonitoringConfig): Promise<void> => http.put('/configs/monitoring/battery', config),
}

export default configAPI
