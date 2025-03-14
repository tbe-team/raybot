import type { RobotState } from '@/types/robot-state'
import type { AxiosRequestConfig } from 'axios'
import http from '@/lib/http'

const robotStateAPI = {
  getRobotState(opt?: Partial<AxiosRequestConfig>): Promise<RobotState> {
    return http.get('/robot-state', opt)
  },
}
export default robotStateAPI
