import type { AxiosRequestConfig } from 'axios'
import type { Alarm, CountAlarmsActivedResponse } from '@/types/alarm'
import type { Paging } from '@/types/paging'
import http from '@/lib/http'

export type AlarmStatus = 'ACTIVE' | 'DEACTIVE'

export interface AlarmParams {
  page?: number
  pageSize?: number
  status?: AlarmStatus
}

const alarmAPI = {
  getCountAlarmsActived: (axiosOpts?: AxiosRequestConfig): Promise<CountAlarmsActivedResponse> => {
    // return http.get('/alarms/count/actived', axiosOpts)
    return Promise.resolve({ count: 5 })
  },

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
