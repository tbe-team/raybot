import type { SystemConfig } from '@/types/system-config'
import http from '@/lib/http'

type SystemConfigParams = SystemConfig

const systemAPI = {
  getSystemConfig(): Promise<SystemConfig> {
    return http.get('/system/config')
  },
  updateSystemConfig(params: SystemConfigParams): Promise<SystemConfig> {
    return http.put('/system/config', params)
  },
  restartSystem(): Promise<void> {
    return http.post('/system/restart')
  },
}
export default systemAPI
