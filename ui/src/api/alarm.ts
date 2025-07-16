import type { AxiosRequestConfig } from 'axios'
import type { Alarm } from '@/types/alarm'
import type { Paging } from '@/types/paging'
import http from '@/lib/http'

export type AlarmStatus = 'ACTIVE' | 'DEACTIVE'

export interface AlarmParams {
  page?: number
  pageSize?: number
  status?: AlarmStatus
}

const alarmAPI = {
  getAlarms: (params: AlarmParams, axiosOpts?: AxiosRequestConfig): Promise<Paging<Alarm>> => {
    return http.get('/alarms', {
      params,
      ...axiosOpts,
    })
  },
  deleteAlarms: (axiosOpts?: AxiosRequestConfig): Promise<void> => {
    return http.delete('/alarms', axiosOpts)
  },
}

export default alarmAPI
