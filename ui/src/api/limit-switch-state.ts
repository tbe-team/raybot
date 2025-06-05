import type { AxiosRequestConfig } from 'axios'
import type { LimitSwitchState } from '@/types/limit-switch-state'
import http from '@/lib/http'

const limitSwitchStateAPI = {
  getLimitSwitchState: (axiosOpts?: Partial<AxiosRequestConfig>): Promise<LimitSwitchState> =>
    http.get('states/limit-switch', axiosOpts),
}

export default limitSwitchStateAPI
