import type { AxiosRequestConfig } from 'axios'
import type { SystemInfo } from '@/types/system-info'
import http from '@/lib/http'

const systemAPI = {
  reboot(): Promise<void> {
    return http.post('/system/reboot')
  },
  stopEmergency(): Promise<void> {
    return http.post('/system/stop-emergency')
  },
  getInfo(axiosOpts?: Partial<AxiosRequestConfig>): Promise<SystemInfo> {
    return http.get('/system/info', axiosOpts)
  },
}
export default systemAPI
