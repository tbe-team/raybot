import type { AxiosRequestConfig } from 'axios'
import type { RobotState } from '@/types/robot-state'
import http from '@/lib/http'

const robotStateAPI = {
  getRobotState(axiosOpts?: Partial<AxiosRequestConfig>): Promise<RobotState> {
    return http.get('/robot-state', axiosOpts)
  },
}
export default robotStateAPI
