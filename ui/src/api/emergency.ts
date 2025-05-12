import type { EmergencyState } from '@/types/emergency'
import type { AxiosRequestConfig } from 'axios'
import http from '@/lib/http'

const emergencyAPI = {
  getEmergencyState: (axiosOpts?: Partial<AxiosRequestConfig>): Promise<EmergencyState> =>
    http.get('/emergency/state', axiosOpts),
  stopEmergency: (axiosOpts?: Partial<AxiosRequestConfig>): Promise<void> =>
    http.post('/emergency/stop', axiosOpts),
  resumeEmergency: (axiosOpts?: Partial<AxiosRequestConfig>): Promise<void> =>
    http.post('/emergency/resume', axiosOpts),
}

export default emergencyAPI
