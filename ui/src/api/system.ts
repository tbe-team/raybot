import type { AxiosRequestConfig } from 'axios'
import type { Info } from '@/types/info'
import http from '@/lib/http'

const systemAPI = {
  reboot(): Promise<void> {
    return http.post('/system/reboot')
  },
  stopEmergency(): Promise<void> {
    return http.post('/system/stop-emergency')
  },
  getInfo(axiosOpts?: Partial<AxiosRequestConfig>): Promise<Info> {
    return http.get('/system/info', { doNotShowLoading: axiosOpts?.doNotShowLoading })
  },
}
export default systemAPI
