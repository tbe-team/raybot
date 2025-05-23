import type { AxiosRequestConfig } from 'axios'
import type { RobotState } from '@/types/robot-state'
import http from '@/lib/http'

const robotStateAPI = {
  getRobotState(opts?: Partial<AxiosRequestConfig>): Promise<RobotState> {
    return http.get('/robot-state', opts)
  },
}
export default robotStateAPI
