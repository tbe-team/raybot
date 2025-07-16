import type { AxiosRequestConfig } from 'axios'
import type { AlarmStatus } from '@/api/alarm'
import { useMutation, useQuery } from '@tanstack/vue-query'
import alarmAPI from '@/api/alarm'

export const COUNT_ALARMS_ACTIVED_QUERY_KEY = 'countAlarmsActived'
export const ALARMS_QUERY_KEY = 'alarms'

export function useCountAlarmsActivedQuery(opts?: { axiosOpts?: Partial<AxiosRequestConfig>, refetchInterval?: number }) {
  return useQuery({
    queryKey: [COUNT_ALARMS_ACTIVED_QUERY_KEY],
    queryFn: () => alarmAPI.getAlarms({ page: 1, pageSize: 1, status: 'ACTIVE' }, opts?.axiosOpts),
    refetchInterval: opts?.refetchInterval,
    select: data => ({ count: data.totalItems }),
  })
}

export function useAlarmsQuery(page: Ref<number>, pageSize: Ref<number>, status: Ref<AlarmStatus>, opts?: { axiosOpts?: Partial<AxiosRequestConfig>, refetchInterval?: number }) {
  return useQuery({
    queryKey: [ALARMS_QUERY_KEY, page, pageSize, status],
    queryFn: () => alarmAPI.getAlarms({
      page: page.value,
      pageSize: pageSize.value,
      status: status.value,
    }, opts?.axiosOpts),
    refetchInterval: opts?.refetchInterval,
  })
}

export function useDeleteAlarmMutation(opts?: { axiosOpts?: Partial<AxiosRequestConfig> }) {
  return useMutation({
    mutationFn: () => alarmAPI.deleteAlarms(opts?.axiosOpts),
  })
}
