import type { RobotState } from '@/types/robot-state'
import http from '@/lib/http'

const robotStateAPI = {
  getRobotState(): Promise<RobotState> {
    return http.get('/robot-state', {
      // @ts-expect-error doNotShowLoading is not in AxiosRequestConfig
      doNotShowLoading: true,
    })
  },
}
export default robotStateAPI
